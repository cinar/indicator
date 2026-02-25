# Momentum Package - Indicators

The `momentum` package provides a collection of indicators for analyzing price momentum and oscillator signals.

## Key Indicators

- **Oscillators:** `AwesomeOscillator`, `ChaikinOscillator`, `StochasticOscillator`, `Rsi` (Relative Strength Index), `Rvi` (Relative Vigor Index), `WilliamsR`.
- **Specialized:** `ConnorsRSI`, `Fisher` (Fisher Transform), `IchimokuCloud`, `Ppo` (Percentage Price Oscillator), `Pvo` (Percentage Volume Oscillator), `Qstick`.
- **Trends:** `PringsSpecialK`, `TdSequential`.

## Implementation Pattern

All momentum indicators are generic over `helper.Number` and follow the `Compute(<-chan T) <-chan T` pattern.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `GEMINI.md`), using `testdata/*.csv` files (e.g., `rsi.csv`, `awesome_oscillator.csv`) for historical data validation.
