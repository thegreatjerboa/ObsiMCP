/* 
package main

import (
	"obsimcp/tools"

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
    noteTools := tools.NewNoteTool()
    s.AddTool(noteTools.ReadNote())
    

    return s
}
*/
