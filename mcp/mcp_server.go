package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cinar/indicator/v2/strategy"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// BacktestArgs represents the arguments for the backtest MCP tool
type BacktestArgs struct {
	Strategy StrategyType `json:"strategy"`
	Data     OhlcvData    `json:"data"`
}

// RunMCPServer starts the MCP server for the backtest functionality
func RunMCPServer() *server.MCPServer {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Backtest MCP server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add backtest tool with schema
	tool := mcp.NewTool("backtest",
		mcp.WithDescription("Run a backtest with the specified strategy and OHLCV data"),
		mcp.WithString("strategy",
			mcp.Required(),
			mcp.Description("The strategy to use for backtesting"),
			mcp.Enum(GetAllStrategyTypes()...),
		),
		mcp.WithObject("data",
			mcp.Required(),
			mcp.Description("OHLCV data for backtesting"),
			mcp.Properties(map[string]any{
				"date": map[string]any{
					"type":        "array",
					"description": "Array of timestamps (Unix seconds)",
					"items":       map[string]any{"type": "integer"},
				},
				"opening": map[string]any{
					"type":        "array",
					"description": "Array of opening prices",
					"items":       map[string]any{"type": "number"},
				},
				"closing": map[string]any{
					"type":        "array",
					"description": "Array of closing prices",
					"items":       map[string]any{"type": "number"},
				},
				"high": map[string]any{
					"type":        "array",
					"description": "Array of high prices",
					"items":       map[string]any{"type": "number"},
				},
				"low": map[string]any{
					"type":        "array",
					"description": "Array of low prices",
					"items":       map[string]any{"type": "number"},
				},
				"volume": map[string]any{
					"type":        "array",
					"description": "Array of volume values",
					"items":       map[string]any{"type": "number"},
				},
			}),
		),
	)

	// Add tool handler using the typed handler
	s.AddTool(tool, mcp.NewTypedToolHandler(handleBacktest))

	return s
}

// handleBacktest processes the backtest request and returns the actions
func handleBacktest(ctx context.Context, request mcp.CallToolRequest, args BacktestArgs) (*mcp.CallToolResult, error) {
	results, err := runBacktest(args.Strategy, args.Data)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process strategy: %v", err)), nil
	}

	if len(results) == 0 || len(results[0].Transactions) == 0 {
		return mcp.NewToolResultError("no transaction data available"), nil
	}

	// Convert strategy.Action to numeric values (1=BUY, 0=HOLD, -1=SELL)
	actions := make([]int, len(results[0].Transactions))
	for i, action := range results[0].Transactions {
		switch action {
		case strategy.Buy:
			actions[i] = 1
		case strategy.Sell:
			actions[i] = -1
			// strategy.Hold is 0 by default
		}
	}

	// Create the response JSON
	response := struct {
		Actions []int   `json:"actions"`
		Outcome float64 `json:"outcome"`
	}{
		Actions: actions,
		Outcome: results[0].Outcome,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// GetAllStrategyTypes returns a slice of all available strategy types as strings
func GetAllStrategyTypes() []string {
	return []string{
		// Base strategies
		string(StrategyBuyAndHold),

		// Trend strategies
		string(StrategyAlligator),
		string(StrategyAroon),
		string(StrategyApo),
		string(StrategyBop),
		string(StrategyCci),
		string(StrategyDema),
		string(StrategyGoldenCross),
		string(StrategyKama),
		string(StrategyKdj),
		string(StrategyMACD),
		string(StrategyQstick),
		string(StrategySmma),
		string(StrategyTrima),
		string(StrategyTripleMaCrossover),
		string(StrategyTsi),
		string(StrategyVwma),
		string(StrategyWeightedClose),

		// Momentum strategies
		string(StrategyAwesomeOscillator),
		string(StrategyRsi),
		string(StrategyStochasticRsi),
		string(StrategyTripleRsi),

		// Volume strategies
		string(StrategyChaikinMoneyFlow),
		string(StrategyEaseOfMovement),
		string(StrategyForceIndex),
		string(StrategyMoneyFlowIndex),
		string(StrategyNegativeVolumeIndex),
		string(StrategyWeightedAveragePrice),
	}
}
