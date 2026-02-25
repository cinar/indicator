# Strategy Package - Trading Logic

The `strategy` package defines the core interfaces and logic for generating buy, sell, and hold signals from indicator data.

## Key Components

- **Models:** `Action`, `Outcome`, `Result`.
- **Interfaces:** `Strategy`.
- **Combinators:** `AndStrategy`, `OrStrategy`, `MajorityStrategy`, `SplitStrategy`.
- **Predefined:** `BuyAndHoldStrategy`.

## Strategy Interface

The `Strategy` interface defines how actions are generated for each data snapshot.
```go
type Strategy[T helper.Number] interface {
	Compute(snapshots <-chan asset.Snapshot) <-chan Action
	IdlePeriod() int
	String() string
}
```

## Strategy Logic

Strategies take a stream of `asset.Snapshot` data and return a stream of `Action` values (Buy, Sell, Hold).
- `AndStrategy`: Requires all sub-strategies to signal the same action.
- `OrStrategy`: Signals the action if any sub-strategy does.
- `MajorityStrategy`: Signals the action if most sub-strategies agree.
- `SplitStrategy`: Divides data into segments for independent analysis.

## Testing Standard

Strategies are tested using historical snapshot data from CSV files (as detailed in the root `GEMINI.md`), ensuring buy/sell/hold signals are correctly generated.
