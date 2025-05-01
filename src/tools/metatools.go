package tools

import (
	"context"
	"fmt"
	"obsimcp/src/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

type MetaTools interface{
    // GetNoteFrontmatter() Get Note Frontmatter Info
    GetNoteFrontmatter() (tool mcp.Tool, handler server.ToolHandlerFunc)
    // AddFrontmatter() Add Frontmatter Info To A Note
    AddFrontmatter() (tool mcp.Tool, handler server.ToolHandlerFunc)
    // GetNoteTags() Get A Note Tags
    GetNoteTags() (tool mcp.Tool, handler server.ToolHandlerFunc)
}

type metaTools struct{

}

func NewMetaTools() MetaTools {
    return &metaTools{}
}

func (m *metaTools) GetNoteFrontmatter() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "GetNoteFrontmatter",
        mcp.WithDescription(`Get the frontmatter YAML metadata from a specific Obsidian note`),
        mcp.WithString(
            "note_path", 
            mcp.Required(), 
            mcp.Description(`Full path to the note`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        note_path := request.Params.Arguments["note_path"].(string)

        if !strings.HasPrefix(note_path, config.Cfg.Vault.Path) {
            return mcp.NewToolResultError("Notepath is a illegalpath"), nil
        }

        if _, err := os.Stat(note_path); os.IsNotExist(err) {
            return mcp.NewToolResultError("File not exist"), nil
        }

        if filepath.Ext(note_path) != ".md" {
            return mcp.NewToolResultError("It's not a markdown"), nil
        }

        content, err := os.ReadFile(note_path)
        if err != nil {
            return mcp.NewToolResultError("Failed to read note"), err
        }

        text := string(content)
        if !strings.HasPrefix(text, "---") {
            return mcp.NewToolResultError("Failed to find frontmatter"), nil
        }

        sections := strings.SplitN(text, "---", 3)
        if len(sections) < 3 {
            return mcp.NewToolResultError("Frontmatter format error"), nil
        }

        frontmatter:= strings.TrimSpace(sections[1])

        return mcp.NewToolResultText(frontmatter), nil
    }

    return
}

func (m *metaTools) AddFrontmatter() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "AddFrontmatter",
        mcp.WithDescription(`Add frontmatter metadata to an Obsidian note. Will fail if frontmatter already exists.`),
        mcp.WithString(
            "note_path", 
            mcp.Required(), 
            mcp.Description(`Full path to the note`),
        ),
        mcp.WithString(
            "frontmatter", 
            mcp.Required(), 
            mcp.Description(`YAML frontmatter content to add, excluding the '---' markers`),
        ),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        notePath := request.Params.Arguments["note_path"].(string)
        frontmatterContent := request.Params.Arguments["frontmatter"].(string)

        if _, err := os.Stat(notePath); err != nil {
            if os.IsNotExist(err) {
                return mcp.NewToolResultError("Note file does not exist"), nil
            }
            return mcp.NewToolResultError("Error accessing file"), err
        }

        if filepath.Ext(notePath) != ".md" {
            return mcp.NewToolResultError("Not a markdown file"), nil
        }

        contentBytes, err := os.ReadFile(notePath)
        if err != nil {
            return mcp.NewToolResultError("Failed to read file"), err
        }

        content := string(contentBytes)

        // check if frontmatter already exists
        if strings.HasPrefix(content, "---") {
            return mcp.NewToolResultError("Frontmatter already exists"), nil
        }

        // prepend frontmatter
        newContent := fmt.Sprintf("---\n%s\n---\n\n%s", frontmatterContent, content)
        err = os.WriteFile(notePath, []byte(newContent), 0644)
        if err != nil {
            return mcp.NewToolResultError("Failed to write new content"), err
        }

        return mcp.NewToolResultText("Frontmatter added successfully"), nil
    }

    return
}

func (m *metaTools) GetNoteTags() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    tool = mcp.NewTool(
        "GetNoteTags",
        mcp.WithDescription("Extract the 'tags' field from the frontmatter of a markdown note and return as plain text."),
        mcp.WithString("note_path", mcp.Required(), mcp.Description("Full path to the markdown note file")),
    )

    handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        notePath := request.Params.Arguments["note_path"].(string)

        // check if file exists
        if _, err := os.Stat(notePath); err != nil {
            if os.IsNotExist(err) {
                return mcp.NewToolResultError("Error: Note file does not exist."), nil
            }
            return mcp.NewToolResultError("Error: Failed to access note file."), err
        }

        // must be a markdown file
        if filepath.Ext(notePath) != ".md" {
            return mcp.NewToolResultError("Error: Not a markdown (.md) file."), nil
        }

        // read file content
        contentBytes, err := os.ReadFile(notePath)
        if err != nil {
            return mcp.NewToolResultError("Error: Failed to read file."), err
        }

        content := string(contentBytes)

        // frontmatter must start with '---'
        if !strings.HasPrefix(content, "---") {
            return mcp.NewToolResultError("Error: No frontmatter found."), nil
        }

        // extract frontmatter lines
        lines := strings.Split(content, "\n")
        var yamlLines []string
        for i := 1; i < len(lines); i++ {
            if lines[i] == "---" {
                break
            }
            yamlLines = append(yamlLines, lines[i])
        }

        yamlContent := strings.Join(yamlLines, "\n")

        // parse YAML
        var meta map[string]interface{}
        err = yaml.Unmarshal([]byte(yamlContent), &meta)
        if err != nil {
            return mcp.NewToolResultError("Error: Failed to parse YAML frontmatter."), err
        }

        // extract tags
        var tags []string
        switch t := meta["tags"].(type) {
        case string:
            tags = append(tags, t)
        case []interface{}:
            for _, item := range t {
                if str, ok := item.(string); ok {
                    tags = append(tags, str)
                }
            }
        }

        if len(tags) == 0 {
            return mcp.NewToolResultText("Tags: none"), nil
        }

        return mcp.NewToolResultText(fmt.Sprintf("Tags: %s", strings.Join(tags, ", "))), nil
    }

    return
}
