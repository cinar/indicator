package main

import (
	"flag"
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	httpMode := flag.Bool("http", false, "Run in HTTP server mode")
	flag.Parse()

	log.Println("Starting MCP server...")
	s := RunMCPServer()
	if *httpMode {
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :8080/mcp")
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
