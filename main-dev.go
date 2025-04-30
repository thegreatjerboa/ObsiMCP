package main

import (
	"fmt"
	"obsimcp/src/config"
	"obsimcp/src/tools"

	"github.com/mark3labs/mcp-go/server"
)

func createServer() *server.MCPServer{
    s := server.NewMCPServer(
        "ObsiMCP",
        "1.0.0",
        server.WithResourceCapabilities(true, true),
        server.WithLogging(),
        server.WithRecovery(),
    )

    // add tools
    // nopte tools
    noteTools := tools.NewNoteTool()
    // s.AddTool(noteTools.ReadNote())
    s.AddTool(noteTools.GetNoteFullPath())
    s.AddTool(noteTools.ReadNoteByFullPath())
    s.AddTool(noteTools.WriteNoteByFullPath())
    s.AddTool(noteTools.CreateANote())
    s.AddTool(noteTools.DeleteNote())
    s.AddTool(noteTools.GetNoteList())
    s.AddTool(noteTools.MoveOneNote())

    // folder tools
    folderTools := tools.NewFolderTools()
    s.AddTool(folderTools.FindAllFolderByName())
    s.AddTool(folderTools.CreateFolder())
    
    // meta tools
    metatools := tools.NewMetaTools()
    s.AddTool(metatools.GetNoteFrontmatter())
    s.AddTool(metatools.AddFrontmatter())
    s.AddTool(metatools.GetNoteTagsText())
    return s
}

func main() {
    // InitConfig
    config.InitConfig()

    s := createServer()

    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server Error: %v", err)
    }
}
