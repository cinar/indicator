# Valuation Package - Indicators

The `valuation` package provides plain time-value-of-money functions, not streaming indicators.

## Key Functions

- `Fv(pv, rate float64, years int) float64`: Future Value of a present value.
- `Pv(fv, rate float64, years int) float64`: Present Value of a future value.
- `Npv(rate float64, cfs []float64) float64`: Net Present Value of a series of cash flows.

## Pattern

Unlike most of the codebase, this package is not channel-based or generic: each function is a plain, synchronous `float64` calculation with no `Compute`/`IdlePeriod`/`String` methods.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `AGENTS.md`), reading `PV`/`rate`/`years`/`FV`-style rows from `testdata/*.csv` and comparing against the rounded function result.
