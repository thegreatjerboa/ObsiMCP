# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ObsiMCP is an MCP (Model Context Protocol) server for automating Obsidian vault operations. It exposes tools that AI assistants can use to read, write, search, and manage notes in an Obsidian vault. Built with Go using the [mcp-go](https://github.com/mark3labs/mcp-go) framework.

## Build and Run

```bash
# Build the executable
go build -o main main.go

# Run tests
go test ./...

# Run a specific test file
go test ./src/plugins/local-rest-api/
```

## Docker Build

```bash
# Build the Docker image
docker build -t obsimcp:local .

# Run with vault mounted
docker run -i --rm \
  -v /path/to/vault:/vault \
  -v /path/to/backup:/backup \
  obsimcp:local
```

Environment variables:
- `VAULT_PATH`: Path inside container for vault (default: `/vault`)
- `BACKUP_PATH`: Path inside container for backups (default: `/backup`)
- `TEMPLATE_PATH`: Path inside container for templates (default: `/templates`)

## Configuration

Edit `src/config/config.yaml` before running:
- `vault.path`: Path to your Obsidian vault
- `backup.path`: Required backup directory for overwrite/delete operations
- `template.path`: Path to templates folder in vault
- `plugins.rest_api`: (Optional) Base URL and auth token for Obsidian Local REST API plugin - only needed for ListTemplates tool

## Architecture

### Entry Point
`main.go` initializes config, the Local REST API client, creates the MCP server, and runs it with stdio transport.

### Server Setup (`src/server.go`)
Creates the MCP server and registers all tool handlers. Tools are organized into four categories:
- **NoteTool**: Note CRUD operations (read, write, create, delete, move, list)
- **FolderTools**: Folder operations (find, create)
- **MetaTools**: Frontmatter/metadata operations (get/add frontmatter, get tags)
- **TemplateTools**: Template listing

### Tool Pattern
Each tool type follows an interface pattern with a factory function:
```go
type NoteTool interface {
    ReadNote() (mcp.Tool, server.ToolHandlerFunc)
    // ...
}
func NewNoteTool() NoteTool { return &noteTool{} }
```

Tools return a tuple of (tool definition, handler function) that gets registered with `s.AddTool()`.

### Local REST API Plugin (`src/plugins/local-rest-api/`)
Client wrapper for Obsidian's Local REST API plugin. Used by template tools to fetch file details and frontmatter via HTTP. Initialized as a global `Client` variable.

### Utils (`src/utils/`)
- `backup.go`: Creates timestamped backups before destructive operations
- `pathcheck.go`: Path validation (vault boundary checks, markdown file checks, existence checks)
- `requests.go`: HTTP request helper with Bearer token auth

### Key Safety Features
- All paths are validated against the vault path to prevent traversal
- Overwrite and delete operations automatically create backups
- Delete only works on `.md` files
