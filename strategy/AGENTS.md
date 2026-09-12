# Strategy Package - Trading Logic

The `strategy` package defines the core interfaces and logic for generating buy, sell, and hold signals from indicator data.

## Key Components

- **Models:** `Action` (Buy/Sell/Hold). `Result` exists only to compare expected vs. actual actions in tests. `Outcome`/`OutcomeWithContext` are functions, not a model — they simulate the P&L of a given action sequence.
- **Metrics:** `SharpeRatioWithContext` computes a risk-adjusted metric from an outcome curve (more may follow, e.g. Sortino, max drawdown).
- **Interface:** `Strategy` (not generic).
- **Combinators:** `AndStrategy`, `OrStrategy`, `MajorityStrategy`, `SplitStrategy`.
- **Predefined:** `BuyAndHoldStrategy`.

## Strategy Interface

The `Strategy` interface defines how actions are generated for each data snapshot.
```go
type Strategy interface {
	Name() string
	Compute(snapshots <-chan *asset.Snapshot) <-chan Action
	Report(snapshots <-chan *asset.Snapshot) *helper.Report
}
```

## Strategy Logic

Strategies take a stream of `*asset.Snapshot` data and return a stream of `Action` values (Buy, Sell, Hold).
- `AndStrategy`: Requires all sub-strategies to signal the same action.
- `OrStrategy`: Signals the action if any sub-strategy does, with no conflicting recommendation from another.
- `MajorityStrategy`: Signals the action most sub-strategies agree on.
- `SplitStrategy`: Uses one strategy to identify Buy opportunities and a separate strategy to identify Sell opportunities; a conflicting recommendation resolves to Hold.

## Testing Standard

Strategies are tested using historical snapshot data from CSV files (as detailed in the root `AGENTS.md`), ensuring buy/sell/hold signals are correctly generated.
