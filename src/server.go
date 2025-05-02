package server

import (
	"obsimcp/src/tools"

	"github.com/mark3labs/mcp-go/server"
)

func CreateServer() *server.MCPServer {
	s := server.NewMCPServer(
		"ObsiMCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// add tools

	// note tools
	noteTools := tools.NewNoteTool()
	// s.AddTool(noteTools.ReadNote())
	s.AddTool(noteTools.GetNote())
	s.AddTool(noteTools.ReadNote())
	s.AddTool(noteTools.WriteNote())
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
	s.AddTool(metatools.GetNoteTags())

	templateTools := tools.NewTemplateTools()
	s.AddTool(templateTools.ListTemplates())

	return s
}

func ServerRunWithStdio(s *server.MCPServer) error {
	if err := server.ServeStdio(s); err != nil {
		return err
	}

	return nil
}
