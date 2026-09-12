# Momentum Package - Indicators

The `momentum` package provides a collection of indicators for analyzing price momentum and oscillator signals.

## Key Indicators

- **Key Indicators:** `AwesomeOscillator`, `ChaikinOscillator`, `StochasticOscillator`, `UltimateOscillator`, `Rsi` (Relative Strength Index), `Rvi` (Relative Vigor Index), `WilliamsR`.
- **Specialized:** `ConnorsRsi`, `Fisher` (Fisher Transform), `IchimokuCloud`, `InternalBarStrength`, `Ppo` (Percentage Price Oscillator), `Pvo` (Percentage Volume Oscillator), `Qstick`.
- **Trends:** `PringsSpecialK`, `TdSequential`.

## Implementation Pattern

Momentum indicators follow the `Compute(<-chan T) <-chan T` pattern and are generic over `helper.Float` (`TdSequential` uses the broader `helper.Number`).

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `AGENTS.md`), using `testdata/*.csv` files (e.g., `rsi.csv`, `awesome_oscillator.csv`) for historical data validation.
