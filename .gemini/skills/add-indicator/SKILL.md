---
name: add-indicator
description: Add a new technical analysis indicator to the project. Use when you need to implement a new indicator (e.g., SMA, EMA, RSI) following the project's streaming data patterns, Go generics, and channel-based API.
---

# Add Indicator

This skill provides the workflow and standards for adding new technical analysis indicators to the `indicator` project.

## Workflow

### 1. Identify Category and Location
Indicators are organized by their technical analysis category. Determine where the new indicator belongs:
- `trend/`: Trend-following indicators (SMA, EMA, MACD, etc.)
- `momentum/`: Momentum indicators (RSI, Stochastic, etc.)
- `volatility/`: Volatility indicators (Bollinger Bands, ATR, etc.)
- `volume/`: Volume-based indicators (OBV, Chaikin Money Flow, etc.)
- `valuation/`: Asset valuation (FV, NPV, etc.)

### 2. Implementation Standards

#### Streaming Data Patterns
- Indicators MUST operate on unlimited data streams using Go routines and channels.
- Use `<-chan T` for inputs and `<-chan T` (or multiple channels) for outputs.
- Indicators may have an "idle" or "warm-up" period where they don't produce output until their internal window is filled.

#### Generics and Types
- All indicators must use Go generics with the `helper.Number` constraint.
- Define a struct for the indicator and a `New<Indicator>` factory function.
- Implement a `Compute` method: `func (i *Indicator[T]) Compute(c <-chan T) <-chan T`.
- Implement an `IdlePeriod` method: `func (i *Indicator[T]) IdlePeriod() int` to return the number of elements it skips.

#### Example Implementation (inspired by `trend/macd.go`)
```go
type MyIndicator[T helper.Number] struct {
    Period int
}

func NewMyIndicator[T helper.Number](period int) *MyIndicator[T] {
    return &MyIndicator[T]{Period: period}
}

func (m *MyIndicator[T]) Compute(c <-chan T) <-chan T {
    // Implementation using helper functions and Go routines
}

func (m *MyIndicator[T]) IdlePeriod() int {
    return m.Period - 1
}
```

### 3. Preventing Deadlocks
- Streaming indicators are prone to deadlocks if not synchronized correctly.
- Use `helper.Skip` (from `helper/skip.go`) to synchronize output streams or handle warm-up periods.
- For repository-level synchronization, refer to `asset/sync.go`.

### 4. Testing Requirements
- **CSV Test Data**: Every indicator MUST have a test data file in a `testdata/` subdirectory (e.g., `trend/testdata/my_indicator.csv`).
- **Data Volume**: Use a large dataset, typically 252 rows (representing a trading year), for consistency.
- **Deadlock Detection**: When running tests for the first time, ALWAYS apply a timeout:
  ```bash
  go test -v ./trend -run TestMyIndicator -timeout 30s
  ```
- **Verification**: Use `helper.CheckEquals` and `helper.RoundDigits` for floating-point comparisons.

### 5. Documentation and Integration
- **Lint and Auto-Docs**: After implementation, run the `task` command. This will execute lint checks and update the module-level `README.md` files automatically.
- **Main README**: Manually update the root `README.md` to add the new indicator to the appropriate list under "Indicators Provided".

## Reference Examples
- **Indicator Code**: See `trend/macd.go` for a complex multi-stream implementation.
- **Test Data**: See `trend/testdata/macd.csv` for the expected CSV format.
- **Helper Utilities**: Explore the `helper/` directory for common operations like `Duplicate`, `Pipe`, `Skip`, etc.
