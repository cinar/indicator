package main

import (
	"context"
	"encoding/json"
	"math"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestMCPServer(t *testing.T) {
	client, err := client.NewInProcessClient(RunMCPServer())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Start the client
	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	testInitialize(t, client)
	testPing(t, client)
	testListTools(t, client)
	testCallTool(t, client)
}

func testInitialize(t *testing.T, client *client.Client) {
	t.Run("Initialize", func(t *testing.T) {
		// Initialize
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "test-client",
			Version: "1.0.0",
		}

		result, err := client.Initialize(context.Background(), initRequest)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		if result.ServerInfo.Name != "Backtest MCP server" {
			t.Errorf(
				"Expected server name 'Backtest MCP server', got '%s'",
				result.ServerInfo.Name,
			)
		}
	})
}

func testPing(t *testing.T, client *client.Client) {
	t.Run("Ping", func(t *testing.T) {
		// Test Ping
		if err := client.Ping(context.Background()); err != nil {
			t.Errorf("Ping failed: %v", err)
		}
	})
}

func testListTools(t *testing.T, client *client.Client) {
	t.Run("ListTools", func(t *testing.T) {
		// Test ListTools
		toolsRequest := mcp.ListToolsRequest{}
		toolListResult, err := client.ListTools(context.Background(), toolsRequest)
		if err != nil {
			t.Errorf("ListTools failed: %v", err)
		}
		if toolListResult == nil || len(toolListResult.Tools) == 0 {
			t.Errorf("Expected one tool")
		}
		if toolListResult.Tools[0].Name != "backtest" {
			t.Errorf("Expected tool name 'backtest'")
		}
	})
}

func testCallTool(t *testing.T, client *client.Client) {
	t.Run("CallTool", func(t *testing.T) {
		request := mcp.CallToolRequest{}
		request.Params.Name = "backtest"
		request.Params.Arguments = map[string]any{
			"strategy": "buy_and_hold",
			"data": map[string]any{
				"date":    []int64{1609459200, 1609545600, 1609632000},
				"opening": []float64{100.0, 101.5, 102.0},
				"closing": []float64{101.0, 101.0, 103.0},
				"high":    []float64{101.5, 102.0, 103.5},
				"low":     []float64{99.5, 100.5, 101.5},
				"volume":  []float64{1000, 1500, 2000},
			},
		}

		result2, err2 := client.CallTool(context.Background(), request)
		if err2 != nil {
			t.Fatalf("CallTool failed: %v", err2)
		}

		if result2 == nil {
			t.Fatalf("Expected a result")
		}

		if result2.IsError {
			t.Fatalf("Expected no error")
		}

		if len(result2.Content) != 1 {
			t.Fatalf("Expected 1 content item, got %d", len(result2.Content))
		}

		textContent, ok := result2.Content[0].(mcp.TextContent)
		if !ok {
			t.Errorf("Expected text content")
		}

		// Parse the JSON content
		var jsonResult struct {
			Actions []int   `json:"actions"`
			Outcome float64 `json:"outcome"`
		}
		if err := json.Unmarshal([]byte(textContent.Text), &jsonResult); err != nil {
			t.Errorf("Failed to parse JSON: %v", err)
		}

		if math.Abs(jsonResult.Outcome-0.019802) > 1e-6 {
			t.Errorf("Expected outcome 0.019802, got %f", jsonResult.Outcome)
		}

		if len(jsonResult.Actions) != 3 {
			t.Errorf("Expected 3 actions, got %d", len(jsonResult.Actions))
		}
	})
}
