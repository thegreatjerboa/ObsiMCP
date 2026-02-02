package main

import (
	"fmt"
	server "obsimcp/src"
	"obsimcp/src/config"
	api "obsimcp/src/plugins/local-rest-api"
)

func main() {
	// InitConfig
	config.InitConfig()

	// Initialize REST API client only if configured
	// (optional - enables template tools via Obsidian Local REST API plugin)
	if config.Cfg.Plugins.Rest.BaseUrl != "" {
		err := api.InitLocalRestApi()
		if err != nil {
			fmt.Printf("Warning: Local REST API not available: %v\n", err)
			// Continue without REST API - basic tools will still work
		}
	}

	// create server
	s := server.CreateServer()

	// Run server with stdio
	if err := server.ServerRunWithStdio(s); err != nil {
		fmt.Printf("Server run failed: %v\n", err)
	}
}
