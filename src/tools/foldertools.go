package tools

import (
	"context"
	"encoding/json"
	"obsimcp/src/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type FolderTools interface{
    // FindAllFolderByName() According Foolder Name Find All Folder
    FindAllFolderByName() (mcp.Tool, server.ToolHandlerFunc)
    // CreateFolder()
    CreateFolder() (mcp.Tool, server.ToolHandlerFunc)
}

type folderTools struct{

}

func NewFolderTools() FolderTools{
    return &folderTools{}
}

func (ft *folderTools) CreateFolder() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "CreateFolder",
        mcp.WithDescription(`Create a folder in Vault based on the full path to the folder`),
        mcp.WithString(
            "folder_path",
            mcp.Required(),
            mcp.Description(`The full path to the folder being created (including the vault path)`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        folderPath, _ := request.Params.Arguments["folder_path"].(string)
        
        if !strings.HasPrefix(folderPath, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("Folder_path is a illegalpath"), nil
        }
        
        if _, err := os.Stat(folderPath); err == nil {
            return mcp.NewToolResultError("Folder alread exist"), nil
        }

        err := os.MkdirAll(folderPath, 0755)
        if err != nil {
            return mcp.NewToolResultError("Failed to create folder"), err
        }

        return mcp.NewToolResultText("Create folder successfully"), nil
    }

    return 
}



func (ft *folderTools) FindAllFolderByName() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "FindAllFolderByName",
        mcp.WithDescription(`According to the folder name specified by the user, 
        search for all folders with the same name in Vault for the user to choose`),
        mcp.WithString(
            "folder_name",
            mcp.Required(),
            mcp.Description(`Specify the name of the folder`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        folder_name, _ := request.Params.Arguments["folder_name"].(string)
        var matches []string

        err := filepath.Walk(config.Cfg.Vault.Path, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if info.IsDir() && info.Name() == folder_name {
                matches = append(matches, path)
            }

            return nil
        })

        if err != nil {
            return mcp.NewToolResultError("Floder Find Error"), err
        }
        
        msg := ""
        if len(matches) > 1 {
            msg = "Multiple folders with the same name are found, and the user needs to select one"
        }

        jsonData, err := json.Marshal(map[string]any{
            "message": msg,
            "matches": matches,
        })

        if err != nil {
            return mcp.NewToolResultError("Json Marshal Error"), err
        }
        

        return mcp.NewToolResultText(string(jsonData)), nil
    }

    return 
}

