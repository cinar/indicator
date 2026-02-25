# Indicator - Project Overview

Indicator is a Golang library for technical analysis, providing a wide range of indicators, strategies, and backtesting capabilities. It leverages Go 1.22+ generics and channels for streaming data processing.

## Core Architecture

- **Streaming:** Most indicators process data through channels (`<-chan T`), allowing for efficient, pipeline-based calculations.
- **Generics:** Indicators are generic over `helper.Number` (integers and floats).
- **Package Organization:**
  - `asset`: Asset data management and repositories.
  - `backtest`: Framework for testing strategies against historical data.
  - `helper`: Low-level channel and math utilities.
  - `strategy`: Logic for generating buy/sell/hold signals.
  - `trend`, `momentum`, `volatility`, `volume`, `valuation`: Category-specific indicators.

## Development Standards

- **No External Dependencies:** This project aims to have no external dependencies. Do not add any new dependencies.
- **Composition & Reusability:** Build and utilize reusable blocks, particularly those in the `helper/` package. Avoid re-implementing existing logic within a single indicator. For example, if an indicator uses a moving average, it should employ the existing implementation rather than duplicating the logic internally.
- **Copyright Header:** Every file must start with the copyright notice:
  ```go
  // Copyright (c) 2021-2026 Onur Cinar.
  // The source code is provided under GNU AGPLv3 License.
  // https://github.com/cinar/indicator
  ```
- **Indicator Pattern:**
  - Struct named after the indicator (e.g., `Ema[T helper.Number]`).
  - `New[Indicator]` and `New[Indicator]With[Param]` constructors.
  - `Compute(<-chan T) <-chan T` for the main logic.
  - `IdlePeriod() int` to indicate when it starts producing values.
  - `String() string` for a descriptive name.
- **Testing:**
  - Packages should use `[package]_test` suffix (e.g., `package trend_test`).
  - Use `helper` utilities for test setup: `SliceToChan`, `ChanToSlice`, `RoundDigits`.
  - 100% coverage is required for all indicators.

## Testing Standards

All indicators and strategies must be tested using historical data from CSV files (typically located in `testdata/`).

### Test Pattern

1.  **Define Data Struct:** Create a struct with tags matching the CSV headers.
2.  **Read CSV:** Use `helper.ReadFromCsvFile` to load the data.
3.  **Prepare Streams:** Use `helper.Duplicate` to branch the input for processing and validation.
4.  **Extract Fields:** Use `helper.Map` to create specific data channels (e.g., closing prices).
5.  **Align Streams:** Use `helper.Skip` on the expected data stream to account for the indicator's `IdlePeriod()`.
6.  **Validate:** Iterate through the channels and compare `actual` vs `expected`.

```go
func TestIndicator(t *testing.T) {
    type Data struct {
        Close    float64 `header:"Close"`
        Expected float64 `header:"Expected"`
    }

    input, err := helper.ReadFromCsvFile[Data]("testdata/indicator.csv")
    if err != nil {
        t.Fatal(err)
    }

    inputs := helper.Duplicate(input, 2)
    closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

    ind := trend.NewIndicator[float64]()
    actuals := ind.Compute(closing)
    actuals = helper.RoundDigits(actuals, 2)

    inputs[1] = helper.Skip(inputs[1], ind.IdlePeriod())

    for data := range inputs[1] {
        actual := <-actuals
        if actual != data.Expected {
            t.Fatalf("actual %v expected %v", actual, data.Expected)
        }
    }
}
```

## Build and Test Commands

- **Run all checks:** `task`
- **Format:** `task fmt`
- **Lint:** `task lint`
- **Test:** `task test`
- **Build CLI Tools:** `task build-tools`
