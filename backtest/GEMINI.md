# Backtest Package - Performance Evaluation

The `backtest` package provides a framework for testing trading strategies against historical price data and evaluating their performance.

## Key Components

- **Evaluator:** `Backtest`, `Run`.
- **Reports:** `Report`, `HtmlReport`, `DataReport`, `DataStrategyResult`.
- **Factories:** `ReportFactory`, `ReportConfig`.

## Strategy Backtesting

The `Backtest` struct coordinates the execution of a `strategy.Strategy` on historical `asset.Snapshot` data.
```go
type Backtest[T helper.Number] struct {
	Strategy strategy.Strategy[T]
	Snapshots <-chan asset.Snapshot
}
```

## Reporting

- `HtmlReport`: Generates detailed visual performance metrics for a backtested strategy.
- `DataReport`: Provides structured raw data from the backtesting results.
- `DataStrategyResult`: Summary of profit/loss, trades, and drawdown.

## Testing Standard

Backtesting results are validated using historical asset data from CSV files (as detailed in the root `GEMINI.md`), ensuring strategy execution aligns with expected trade outcomes.
