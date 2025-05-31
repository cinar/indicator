# Market Condition Predictor (MCP)

MCP is a Go-based backtesting framework for financial market strategies. It provides a flexible way to test trading strategies against historical OHLCV (Open, High, Low, Close, Volume) market data.

## Features

- Multiple built-in trading strategies
- Support for standard OHLCV data format
- HTTP server mode for remote access
- Standard I/O mode for local usage
- Extensible architecture for adding new strategies

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/indicator.git
   cd indicator/mcp
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o backtest
   ```

## Usage

### Running in HTTP Server Mode

To start the MCP server with HTTP interface:

```bash
./backtest -http
```

The server will start on `http://localhost:8080/mcp`

### Running in Standard I/O Mode

To run in standard I/O mode (for local usage or piping):

```bash
./backtest
```

### API Usage

#### Backtest Endpoint

Make a POST request to `/mcp` with the following JSON payload:

```json
{
  "tool": "backtest",
  "args": {
    "strategy": "rsi",  // see `strategy_factory.go` for available strategies
    "data": {
      "date": [1609459200, 1609545600, 1609632000],
      "opening": [100.0, 101.5, 102.0],
      "closing": [101.0, 101.0, 103.0],
      "high": [101.5, 102.0, 103.5],
      "low": [99.5, 100.5, 101.5],
      "volume": [1000, 1500, 2000]
    }
  }
}
```

#### Available Strategies

See `strategy_factory.go` for available strategies.

## Input Data Format

The input data should be in OHLCV format with the following structure:

```json
{
  "date": [1609459200, 1609545600, 1609632000],
  "opening": [100.0, 101.5, 102.0],
  "closing": [101.0, 101.0, 103.0],
  "high": [101.5, 102.0, 103.5],
  "low": [99.5, 100.5, 101.5],
  "volume": [1000, 1500, 2000]
}
```

## Response Format

The API returns an array of actions with the following structure:

```json
{
  "actions": [1, 0, 0]
}
```

Actions are encoded as:

- 1: BUY
- 0: HOLD
- -1: SELL

### Testing

Run MCP inspector to test the backtest endpoint:

```bash
npx @modelcontextprotocol/inspector
```
