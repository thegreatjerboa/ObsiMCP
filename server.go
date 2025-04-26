package main

import "github.com/mark3labs/mcp-go/server"


func createServer() {
    s := server.NewMCPServer(
        "ObsiMCP",
        "1.0.0",
        server.WithResourceCapabilities(true, true),
        server.WithLogging(),
        server.WithRecovery(),
    )
    // add tools
    


}
