package main

import (
	"fmt"
	"obsimcp/config"
	"obsimcp/tools"

	// "os"

	"github.com/mark3labs/mcp-go/server"
	// "golang.org/x/mod/sumdb/note"
	// "go.uber.org/zap"
)

/*
var logger *zap.Logger

func initLogger() {
    var err error
    cfg := zap.NewProductionConfig()
    cfg.OutputPaths = []string{
		"server.log", // 日志文件
		"stdout",     // 同时输出到控制台
	}
	logger, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
}
*/

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
    s.AddTool(noteTools.GetNoteFullPath())
    s.AddTool(noteTools.ReadNoteByFullPath())
    s.AddTool(noteTools.WriteNoteByFullPath())
    s.AddTool(noteTools.CreateANote())
    

    return s
}

func main() {
    // initLogger()
	// defer logger.Sync()
    // logger.Info("Starting ObsiMCP server...")

    // InitConfig
    config.InitConfig()

    s := createServer()

    if err := server.ServeStdio(s); err != nil {
        //logger.Error("Server error", zap.Error(err))
        // os.Exit(-1)
        fmt.Printf("Server Error: %v", err)
    }
}
