package main

import (
	"flag"
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	transport := flag.String("transport", "studio", "Transport mode: 'http' or 'studio'")
	flag.Parse()

	log.Println("Starting MCP server...")
	s := RunMCPServer()
	switch *transport {
	case "http":
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :8080/mcp")
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "studio":
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf("Invalid transport: %s", *transport)
	}
}
