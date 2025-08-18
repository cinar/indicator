package main

import (
	"flag"
	"log"

	"github.com/mark3labs/mcp-go/server"
)

// main is the entry point of the Backtest MCP server.
// It initializes and starts the server in one of the available transport modes.
//
// Command-line flags:
//
//	-transport string
//	  	Transport mode to use (default "stdio")
//	  	Options:
//	  	  - http:   Run as an HTTP server on port 8080
//	  	  - stdio: Run in standard I/O mode for IDE integration
//
// Examples:
//   - Run in HTTP mode:    ./backtest -transport=http
//   - Run in stdio mode:  ./backtest -transport=stdio
//   - Run with default:    ./backtest
//
// The server provides backtesting functionality for trading strategies
// through the specified transport layer.
func main() {
	transport := flag.String("transport", "stdio", "Transport mode: 'http' or 'stdio'")
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
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf("Invalid transport: %s", *transport)
	}
}
