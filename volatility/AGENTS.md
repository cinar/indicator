# Volatility Package - Indicators

The `volatility` package provides a collection of indicators for analyzing price volatility and market range.

## Key Indicators

- **Bands & Channels:** `AccelerationBands`, `BollingerBands`, `DonchianChannel`, `KeltnerChannel`.
- **Indicators:** `Atr` (Average True Range), `BollingerBandWidth`, `ChandelierExit`, `MovingStd` (Moving Standard Deviation), `PercentB`, `TrueRange` (True Range).
- **Oscillators:** `Po` (Price Oscillator), `SuperTrend`, `UlcerIndex`.

## Pattern

Most volatility indicators return a single channel like `Atr`, but a few return multiple channels for their component bands/lines (e.g., `BollingerBands.ComputeWithContext` returns three channels: upper, middle, lower).

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `AGENTS.md`), using `testdata/*.csv` files (e.g., `atr.csv`, `bollinger_bands.csv`) for historical data validation.
