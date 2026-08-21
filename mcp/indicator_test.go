package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestRunIndicator(t *testing.T) {
	data := OhlcvData{
		Date:    []int64{1609459200, 1609545600, 1609632000, 1609718400, 1609804800, 1609891200, 1609977600},
		Opening: []float64{100, 101, 102, 103, 104, 105, 106},
		Closing: []float64{101, 102, 103, 104, 105, 106, 107},
		High:    []float64{101.5, 102.5, 103.5, 104.5, 105.5, 106.5, 107.5},
		Low:     []float64{99.5, 100.5, 101.5, 102.5, 103.5, 104.5, 105.5},
		Volume:  []float64{1000, 1000, 1000, 1000, 1000, 1000, 1000},
	}

	response, err := runIndicator(context.Background(), IndicatorRequest{
		Indicator: "sma",
		Params:    []string{"Period=3"},
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(response.Dates) != 5 {
		t.Fatalf("expected 5 dates after the idle period, got %d", len(response.Dates))
	}

	sma, ok := response.Series["sma"]
	if !ok {
		t.Fatal(`expected a "sma" series`)
	}

	if len(sma) != len(response.Dates) {
		t.Fatalf("expected %d sma values, got %d", len(response.Dates), len(sma))
	}

	// closings[0:3] = 101,102,103 -> average 102
	if sma[0] != 102 {
		t.Errorf("expected first sma value to be 102, got %v", sma[0])
	}
}

func TestRunIndicatorUnknownIndicator(t *testing.T) {
	_, err := runIndicator(context.Background(), IndicatorRequest{
		Indicator: "does-not-exist",
		Data: OhlcvData{
			Date:    []int64{1609459200},
			Opening: []float64{100},
			Closing: []float64{101},
			High:    []float64{101.5},
			Low:     []float64{99.5},
			Volume:  []float64{1000},
		},
	})
	if err == nil {
		t.Fatal("expected an error for an unknown indicator")
	}
}

func TestRunIndicatorMismatchedLengths(t *testing.T) {
	_, err := runIndicator(context.Background(), IndicatorRequest{
		Indicator: "sma",
		Data: OhlcvData{
			Date:    []int64{1609459200, 1609545600},
			Opening: []float64{100},
			Closing: []float64{101, 102},
			High:    []float64{101.5, 102.5},
			Low:     []float64{99.5, 100.5},
			Volume:  []float64{1000, 1000},
		},
	})
	if err == nil {
		t.Fatal("expected an error for mismatched OHLCV lengths")
	}
}

func TestRunIndicatorInvalidParam(t *testing.T) {
	_, err := runIndicator(context.Background(), IndicatorRequest{
		Indicator: "sma",
		Params:    []string{"Period"},
		Data: OhlcvData{
			Date:    []int64{1609459200},
			Opening: []float64{100},
			Closing: []float64{101},
			High:    []float64{101.5},
			Low:     []float64{99.5},
			Volume:  []float64{1000},
		},
	})
	if err == nil {
		t.Fatal(`expected an error for a param not in "Name=Value" form`)
	}
}

func TestMCPServerIndicatorTool(t *testing.T) {
	c, err := client.NewInProcessClient(RunMCPServer())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	if err := c.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{Name: "test-client", Version: "1.0.0"}

	if _, err := c.Initialize(context.Background(), initRequest); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = "indicator"
	request.Params.Arguments = map[string]any{
		"indicator": "sma",
		"params":    []string{"Period=3"},
		"data": map[string]any{
			"date":    []int64{1609459200, 1609545600, 1609632000, 1609718400, 1609804800},
			"opening": []float64{100, 101, 102, 103, 104},
			"closing": []float64{101, 102, 103, 104, 105},
			"high":    []float64{101.5, 102.5, 103.5, 104.5, 105.5},
			"low":     []float64{99.5, 100.5, 101.5, 102.5, 103.5},
			"volume":  []float64{1000, 1000, 1000, 1000, 1000},
		},
	}

	result, err := c.CallTool(context.Background(), request)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if result.IsError {
		t.Fatalf("Expected no error, got: %+v", result.Content)
	}

	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Content))
	}

	if _, ok := result.Content[0].(mcp.TextContent); !ok {
		t.Errorf("Expected text content")
	}
}
