package tools

import (
	"context"
	"fmt"
	"obsimcp/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NoteTool interface{
    // ReadNote() Get Node Content
    ReadNote() (mcp.Tool, server.ToolHandlerFunc)
}

type noteTool struct{

}

func NewNoteTool() NoteTool{
    return &noteTool{}
}

func (n *noteTool) ReadNote() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "ReadNote",
        mcp.WithDescription("Read content from a obsidian markdown file"),
        mcp.WithString("file_path",
            mcp.Required(),
            mcp.Description("Relative path to the file under the vault(e.g., 'subfolder/note.md')"),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        relativePath := request.Params.Arguments["file_path"].(string)
        // VaultPath is set in config
        fullPath := filepath.Join(config.Cfg.Vault.Path, relativePath)

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
