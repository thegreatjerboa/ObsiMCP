package main

import (
	"fmt"
	server "obsimcp/src"
	"obsimcp/src/config"
)

func main() {
    // InitConfig
    config.InitConfig()
    
    // create server
    s := server.CreateServer()
    
    // Run
    if err := server.ServerRunWithStdio(s); err != nil {
        fmt.Printf("Server run failed: %v\n", err)
    }
}

