# Backtest Package - Performance Evaluation

The `backtest` package provides a framework for testing trading strategies against historical price data and evaluating their performance.

## Key Components

- **Evaluator:** `Backtest` (not generic), which runs one or more `strategy.Strategy` values against one or more assets from an `asset.Repository` and writes the results to a `Report`.
- **Reports:** `Report` (interface), `HTMLReport`, `DataReport`, `DataStrategyResult`.
- **Factory:** `NewReport(name, config string) (Report, error)`, backed by `RegisterReportBuilder`.

## Strategy Backtesting

The `Backtest` struct coordinates running strategies against a repository's assets and reporting the results.
```go
type Backtest struct {
	Names      []string
	Strategies []strategy.Strategy
	Workers    int
	LastDays   int
	Logger     *slog.Logger
}
```

## Reporting

- `HTMLReport`: Generates an HTML report with charts for each asset/strategy pairing.
- `DataReport`: Provides structured, in-memory access to the backtest results (`map[string][]*DataStrategyResult`) for programmatic use.
- `DataStrategyResult`: Per asset/strategy result — `Asset`, `Strategy`, `Outcome` (final cumulative return), `Action` (final recommendation), and `Transactions` (all recommendations). No risk-adjusted metrics (Sharpe Ratio aside) or drawdown yet.

## Testing Standard

Backtesting results are validated using historical asset data from CSV files (as detailed in the root `AGENTS.md`), ensuring strategy execution aligns with expected trade outcomes.
