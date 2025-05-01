# ObsiMCP

A lightweight, extendable MCP server for automating operations on Obsidian vaults, powered by [mcp-go](https://github.com/mark3labs/mcp-go).
ObsiMCP helps you efficiently manage your Obsidian notes and build your own personal knowledge base.


## ✨ Features

- 📖 **Read a Markdown note**
- 🔍 **Search all notes with the same file name** across the entire vault.
- ✍️ **Write content to an existing note**, supporting both:
  - `append` mode (adds to the end)
  - `overwrite` mode (replaces the file content)
- 🆕 **Create a new note**
- 🗑️ **Delete a note**
- 📁 **List files and subfolders** within a given folder (non-recursive).
- 🏷 Modify, add, and get the frontmatter in Note
- ...

## 📦 Requirements

- Go 1.20+ (The lower version may work, but I haven't tried it yet.)
- An initialized Obsidian vault directory

## 🚀 Getting Started

```bash
git clone https://github.com/IAMLEIzZ/ObsiMCP.git
cd obsidian-mcp-server
go build -o main main.go
```

Configure your `vault` path via the config/config.yaml file or environment variable as needed.

```
vault:
  path: "your osbidian vault path"
```

### ⚠️ The BackUpDir Config

**To protect your files, you must set up a backup folder. When you use mcp-server to overwrite or delete, a backup file will be automatically generated to prevent your AI assistant from unstable operation.** (You can manually delete outdated backup files regularly, and we will consider adding an automatic deletion function later.)

```
backup:
  path: "your backup dir"
```

## 🧠 Use Cases

- Build agents or copilots that understand and modify your notes.
- Integrate with Obsidian from external apps or LLMs.
- Automate note maintenance, renaming, or content updates.

For example, in deepchat you can start it like this (Please open the file read and write permissions at the same time)
```
{
  "mcpServers": {
    "ObsiMCP": {
      "command": "Your The absolute path to your main executable", 
      "args": [],
      "shell": false,
      "env": {}
    }
  }
}
```

## 🛠 Example Tool Definitions

Each tool in the MCP server corresponds to an operation:

| Tool Name                   |                  Description                       |
|-----------------------------|----------------------------------------------------|
| `ReadNote`          | Reads content from a given Markdown file.          |
| `GetNote`             | Finds all notes in the vault with a matching name. |
| `WriteNote`         | Appends or overwrites content in an existing note. |
| `CreateANote`                 | Creates a new note without initial content.        |
| `DeleteNote`                  | Deletes a note by its name and backup it.                        |
| `GetNoteList`                 | Lists files and subfolders in a given folder.      |
| `MoveOneNote`                 | Move a note to another folder.                     |
| `GetNoteList`                 | Lists files and subfolders in a given folder.      |
| `FindAllFolderByName`         | Find all folders with the same name in the vault according to the folder name.      |
| `CreateFolder`                 | Delete a folder with folder name.      |
| `GetNoteFrontmatter`                | Get a note frontmatter information      |
| `AddFrontmatter`                | Add frontmatter information to a note     |
| `GetNotetags`                |  Add tags to a note     |
| ...                |  ...     |

## 🔗 References
- [mcp-go](https://github.com/mark3labs/mcp-go)


## 📄 License

```
MIT License. Feel free to use, modify, and contribute!
```
