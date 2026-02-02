package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"obsimcp/src/config"
	api "obsimcp/src/plugins/local-rest-api"
	"obsimcp/src/utils"
)

type TemplateTools interface {
	// ListTemplates List all templates
	ListTemplates() (mcp.Tool, server.ToolHandlerFunc)
}

type templateTools struct{}

func NewTemplateTools() TemplateTools {
	return &templateTools{}
}

type templateResponse struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func (t templateTools) ListTemplates() (tool mcp.Tool, handler server.ToolHandlerFunc) {

	tool = mcp.NewTool(
		"ListTemplates",
		mcp.WithDescription(`List all templates in obsidian template folder.\
								it return template's name and description`),
	)
	handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check if REST API client is available
		if api.Client == nil {
			return mcp.NewToolResultError("Template listing requires Obsidian Local REST API plugin - not configured"), nil
		}
		items, err := api.Client.ListDirectory(config.Cfg.Template.Path)
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("Error listing directory: %v", err),
			), err
		}
		var templates []templateResponse
		for _, item := range items {
			// skip directory
			if !utils.CheckIsMd(item) {
				continue
			}
			frontmatter, e := api.Client.GetVaultFileFrontmatter(fmt.Sprintf("%s/%s", config.Cfg.Template.Path, item))
			if e != nil {
				return mcp.NewToolResultError(
					fmt.Sprintf("Error getting file frontmatter: %v, filename: %s", e, item),
				), e
			}
			var desc string
			if _, ok := frontmatter["description"]; ok {
				desc = frontmatter["description"].(string)
			}

			template := templateResponse{
				Name:        item,
				Path:        fmt.Sprintf("%s/%s", config.Cfg.Template.Path, item),
				Description: desc,
			}
			templates = append(templates, template)
		}
		jsonData, err := json.Marshal(templates)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshalling JSON: %v", err)), err
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
	return
}
