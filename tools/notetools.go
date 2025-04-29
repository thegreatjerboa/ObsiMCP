package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"obsimcp/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NoteTool interface{
    // eadNote() Get Note Content
    ReadNote() (mcp.Tool, server.ToolHandlerFunc)
    // ReadNoteByFullPath() Get Note By fullPath
    ReadNoteByFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // GetNoteFullPath() Get Note Full Path by FileName
    GetNoteFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // WriteNoteByFullPath() Write A Note By File Full Path(include vault path)
    WriteNoteByFullPath() (mcp.Tool, server.ToolHandlerFunc)
    // CreateANote() Create A Note By FullPath
    CreateANote() (mcp.Tool, server.ToolHandlerFunc)
}

type noteTool struct{

}

func NewNoteTool() NoteTool{
    return &noteTool{}
}

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

// Read Note By fullPath
func (n *noteTool) ReadNoteByFullPath() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "ReadNoteByFullPath",
        mcp.WithDescription(`Read content from a obsidian markdown file by fullPath.
                            -file_full_path: Note Full Path`),
        mcp.WithString(
            "file_full_path",
            mcp.Required(),
            mcp.Description(`The full file path including vaultpath（e.g. vaultpath/file_name)`),
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
                            - file_name: The note file name specified by the user`),
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
        mcp.WithDescription(`Write content to the Note according to the full path of the Note (including the Vault path). 
                            The writing method is append or overwrite. The default is append.
                            -file_full_path: The full path to the file to be written`),
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
            mcp.Description(`Write mode: append (append) or overwrite (overwrite), the default is append`),
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
        case "append":
            err := os.WriteFile(file_full_path, []byte(content), 0644)
            if err != nil {
                return mcp.NewToolResultError("Failed to write note"), err
            }
        case "overwrite":
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
                            -target_file_path: The full path of the Note being created`),
        mcp.WithString(
            "target_file_path",
            mcp.Required(),
            mcp.Description(`The full path(including the vault path) to the file to be created (e.g. /vault/abc/def/xxx.md)`),
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
