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
	err := api.InitLocalRestApi()
	if err != nil {
		fmt.Printf("Error initializing local REST API: %v\n", err)
		return
	}

	// create server
	s := server.CreateServer()

	// Run server with stdio
	if err := server.ServerRunWithStdio(s); err != nil {
		fmt.Printf("Server run failed: %v\n", err)
	}
}
