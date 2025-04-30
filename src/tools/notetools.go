package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"obsimcp/src/config"
	"obsimcp/src/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NoteTool interface{
    // eadNote() Get Note Content
    // ReadNote() (mcp.Tool, server.ToolHandlerFunc)
    // ReadNoteByFullPath() Get Note By FullPath
    ReadNoteByFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // GetNoteFullPath() Get Note Full Path By FileName
    GetNoteFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // WriteNoteByFullPath() Write A Note By File Full Path(include vault path)
    WriteNoteByFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // CreateANote() Create A Note By FullPath
    CreateANote() (mcp.Tool, server.ToolHandlerFunc)
    // DeleteNote() Delete A Note By FullPath
    DeleteNote() (mcp.Tool, server.ToolHandlerFunc)
    // GetNoteList() Get All Notes And Folders Under A Folder (non-recursive)
    GetNoteList() (mcp.Tool, server.ToolHandlerFunc)
    // MoveNote() Move One Note To Rarget Path
    MoveOneNote() (mcp.Tool, server.ToolHandlerFunc)
}

type noteTool struct{

}

func NewNoteTool() NoteTool{
    return &noteTool{}
}

/*
// Read the contents of a file based on its relative path (relative to the Obsidian Vault Path)
func (n *noteTool) ReadNote() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "ReadNote",
        mcp.WithDescription(`Read content from a obsidian markdown file`),
        mcp.WithString(
            "file_path",
            mcp.Required(),
            mcp.Description(`Relative path to the file under the vault(e.g. 'subfolder/note.md')`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        relativePath := request.Params.Arguments["file_path"].(string)
        // VaultPath is set in config
        fullPath := filepath.Join(config.Cfg.Vault.Path, relativePath)
        // fullPath := filepath.Join("/Users/iamleizz/study/PersonalKnowledgeBase", relativePath)
        // Make sure the final path is still under vault
        if !strings.HasPrefix(fullPath, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("Access denied: invalid path"), nil
        }

        // Check if a file exists
        fi, err := os.Stat(fullPath)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error accessing file: %v", err)), nil
        }
        if fi.IsDir() {
            return mcp.NewToolResultError("The path is a directory, not a file"), nil
        }
        
        // Read file
        data, err := os.ReadFile(fullPath)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), nil
        }
    
        return mcp.NewToolResultText(string(data)), nil
    }

    return
}
*/

// Read Note By fullPath
func (n *noteTool) ReadNoteByFullPath() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "ReadNoteByFullPath",
        mcp.WithDescription(`Read content from a obsidian markdown file by fullPath.
                            -file_full_path: Note Full Path`),
        mcp.WithString(
            "file_full_path",
            mcp.Required(),
            mcp.Description(`The full file path`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Check is a file exists
        fullPath := request.Params.Arguments["file_full_path"].(string)
        fi, err := os.Stat(fullPath)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error accessing file: %v", err)), err
        }
        if fi.IsDir() {
            return mcp.NewToolResultError("The path is a directory, not a file"), err
        }

        // Read file
        data, err := os.ReadFile(fullPath)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), err
        }

        return mcp.NewToolResultText(string(data)), nil
    }

    return
}

// Get the full path of a file by its file name
func (n *noteTool) GetNoteFullPath()(tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "GetNoteFullPath",
        mcp.WithDescription(`According to the Note file name provided by the user, find all files named with the file and the corresponding 
                            path in the Obsidian Note Library for the user to select.
                            -file_name: The note file name specified by the user`),
        mcp.WithString(
            "file_name",
            mcp.Required(),
            mcp.Description(`The file name to search for (without the .md extension)`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        name, _ := request.Params.Arguments["file_name"].(string)
        var matches []string


        err := filepath.Walk(config.Cfg.Vault.Path, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            if !info.IsDir() && strings.TrimSuffix(info.Name(), ".md") == name {
                matches = append(matches, path)
            }

            return nil
        })

        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error Get Path: %v", err)), err
        }
        
        jsonData, err := json.Marshal(map[string]interface{}{
            "matches": matches,
        })
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Error Get Path: %v", err)), err
        }

        return mcp.NewToolResultText(string(jsonData)), nil
    }

    return
}

// Write the content according to the full path of the Note
func (n *noteTool) WriteNoteByFullPath() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    // 传入三个参数，完整地址，内容以及写的方式（追加写 or 覆盖写）
    tool = mcp.NewTool(
        "WriteNoteByFullPath",
        mcp.WithDescription(`Write content to the Note according to the full path of the Note. 
                            The writing method is append or overwrite. The default is append.
                            -file_full_path: The full path to the file to be written
                            -content: What needs to be written
                            -mode: Write mode: append or overwrite. When the user does not explicitly specify to overwrite, 
                            append is used by default.`),
        mcp.WithString(
            "file_full_path",
            mcp.Required(),
            mcp.Description(`First, you need to search for the corresponding complete file path according to the file name, 
            and then write the content according to the complete file path`),
        ),
        mcp.WithString(
            "content",
            mcp.Required(),
            mcp.Description(`What needs to be written`),
        ),
        mcp.WithString(
            "mode",
            mcp.Description(`Write mode: append or overwrite. 
                            When the user does not explicitly specify to overwrite, append is used by default.`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        file_full_path, ok := request.Params.Arguments["file_full_path"].(string)
        if !ok || file_full_path == "" {
            return mcp.NewToolResultError("Param file_full_path is nil"), nil
        }

        content, ok := request.Params.Arguments["content"].(string)
        if !ok || content == "" {
            return mcp.NewToolResultError("Param content is nil"), nil
        }

        mode, _ := request.Params.Arguments["mode"].(string)
        if mode == "" {
            // Default is append
            mode = "append"
        }
        
        // check path
        if _, err := os.Stat(file_full_path); os.IsNotExist(err) {
            // This server is written based on the existing file
            return mcp.NewToolResultError("file not exist"), err
        }

        switch mode{
        case "overwrite":
            // if overwrite, need backup
            _, err := utils.Backupfile(file_full_path)
            if err != nil {
                return mcp.NewToolResultError("Before overwrite, failed to backup this file"), err
            }
            err = os.WriteFile(file_full_path, []byte(content), 0644)
            if err != nil {
                return mcp.NewToolResultError("Failed to write note"), err
            }
        case "append":
            f, err := os.OpenFile(file_full_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                return mcp.NewToolResultError("Failed to open note"), err
            }
            defer f.Close()

            if _, err := f.WriteString("\n" + content); err != nil {
                return mcp.NewToolResultError("Failed to write note"), err
            }
        default:
            return mcp.NewToolResultError("Invalid write mode, should be append or overwrite"), nil
        }

        return mcp.NewToolResultText("The content has been successfully written to the note"), nil
    }
    return 
}

// Create a note
func (n *noteTool) CreateANote() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "CrateANote",
        mcp.WithDescription(`Create a Note according to the full path. 
                            Note that this server only creates a Note and does not write content to it. 
                            To write content, you need to call the server that writes content.
                            Note: Before calling this method, if the Note being created is not under the Vault root directory, 
                            you need to call other methods to find all directories with the same name as the directory 
                            where the Note is located, and return them to the user for selection 
                            (if there is only one directory, it can be created directly).
                            -target_file_path: The full path of the Note being created.`),
        mcp.WithString(
            "target_file_path",
            mcp.Required(),
            mcp.Description(`The full path to the file to be created (e.g. /vault/abc/def/xxx.md)`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        target_file_path := request.Params.Arguments["target_file_path"].(string)
        
        // check IllegalPath
        if !filepath.HasPrefix(target_file_path, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("IllegalPath"), nil
        }
        
        // check path exist
        if _, err := os.Stat(target_file_path); err == nil {
            return mcp.NewToolResultError("file already exist"), err
        }

        // CreateFile
        f, err := os.Create(target_file_path)
        if err != nil {
            return mcp.NewToolResultError("Failed to create note"), err
        }
        f.Close()

        return mcp.NewToolResultText("Note was successfully created"), nil
    }

    return
}

// DeleteNote
func (n *noteTool) DeleteNote() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "DeleteNote",
        mcp.WithDescription(`Given a full path to a Note, delete the Note. Note: This is a dangerous operation, 
                            please confirm the file path to be deleted with the user before calling this server.`),
        mcp.WithString(
            "target_file_path",
            mcp.Required(),
            mcp.Description(`The full path to the file to be delete (e.g. /vault/abc/def/xxx.md)`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        target_file_path := request.Params.Arguments["target_file_path"].(string)

        // check IllegalPath
        if !strings.HasPrefix(target_file_path, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("IllegalPath"), nil
        }
        
        // check is .md
        if !strings.HasSuffix(target_file_path, ".md") {
            return mcp.NewToolResultError("Not a markdown file"), nil
        }

        // check path exist
        if _, err := os.Stat(target_file_path); os.IsNotExist(err) {
            return mcp.NewToolResultError("File not exist"), nil
        }
        
        // before delete note, need backup file
        if _, err := utils.Backupfile(target_file_path); err != nil {
            return mcp.NewToolResultError("Before delete file, failed to backup file"), err
        }
        

        if err := os.Remove(target_file_path); err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to delete note: %v", err)), nil
        }

        return mcp.NewToolResultText(fmt.Sprintf("Note '%s' delete successfully", target_file_path)), nil
    }

    return
}

func (n *noteTool) GetNoteList() (tool mcp.Tool, handler server.ToolHandlerFunc) {

    tool = mcp.NewTool(
        "GetNoteList",
        mcp.WithDescription(`Given a full path to a folder, get the note file name and folder name in the folder and return it.`),
        mcp.WithString(
            "folder_path",
            mcp.Required(),
            mcp.Description(`The full path to the folder where all notes need to be listed`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        folder_path := request.Params.Arguments["folder_path"].(string)

        // check IllegalPath
        if !strings.HasPrefix(folder_path, config.Cfg.Vault.Path) {
                return mcp.NewToolResultError("IllegalPath"), nil
        }
        
        fi, err := os.Stat(folder_path)
        if err != nil {
            return mcp.NewToolResultError("Folder not found"), nil
        }

        if !fi.IsDir() {
            return mcp.NewToolResultError("Path is not a folder"), nil
        }
        
        entries, err := os.ReadDir(folder_path)
        if err != nil {
            return mcp.NewToolResultError("Failed to read folder"), err
        }
        var result []map[string]interface{}

        for _, entry := range entries {
            item := map[string]interface{}{
                "name":    entry.Name(),
                "is_dir":  entry.IsDir(),
            }
            result = append(result, item)
        }
        
        jsondata, err := json.Marshal(map[string]interface{}{
            "entries": result,
        })

        if err != nil {
            return mcp.NewToolResultError("Json Marshal Failed"), err
        }

        return mcp.NewToolResultText(string(jsondata)), nil
    }
    
    return
}


func (n *noteTool) MoveOneNote() (tool mcp.Tool, hander server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "MoveOneNote",
        mcp.WithDescription(`Move a markdown note file from source path to target path`),
        mcp.WithString(
            "source_path",
            mcp.Required(),
            mcp.Description(`The full path of the note to move`),
        ),
        mcp.WithString(
            "target_path",
            mcp.Required(),
            mcp.Description(`The full path to move the note to`),
        ),
    )

    hander = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        source_path, _ := request.Params.Arguments["source_path"].(string)
        target_path, _ := request.Params.Arguments["target_path"].(string)
        
        // check is illegalPath
        if !strings.HasPrefix(source_path, config.Cfg.Vault.Path) || !strings.HasPrefix(target_path, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("Sourcepath or targetpath is illegalPath"), nil
        }

        // check is .md
		if !strings.HasSuffix(source_path, ".md") || !strings.HasSuffix(target_path, ".md") {
			return mcp.NewToolResultError("Only .md files can be moved"), nil
		}

        // check is exist
		if _, err := os.Stat(source_path); os.IsNotExist(err) {
			return mcp.NewToolResultError("Source file does not exist"), nil
		}

        // check target path is exist（if not exist, create it）
		targetDir := filepath.Dir(target_path)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return mcp.NewToolResultError("Failed to create target directory"), err
		}
        
        if err := os.Rename(source_path, target_path); err != nil {
            return mcp.NewToolResultError("Failed to move note"), err
        }

        return mcp.NewToolResultText(fmt.Sprintf("Move note from %s to %s", source_path, target_path)), nil
    }

    return 
}

