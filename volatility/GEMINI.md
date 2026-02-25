# Volatility Package - Indicators

The `volatility` package provides a collection of indicators for analyzing price volatility and market range.

## Key Indicators

- **Bands & Channels:** `AccelerationBands`, `BollingerBands`, `DonchianChannel`, `KeltnerChannel`.
- **Indicators:** `Atr` (Average True Range), `BollingerBandWidth`, `ChandelierExit`, `MovingStd` (Moving Standard Deviation), `PercentB`.
- **Oscillators:** `Po` (Price Oscillator), `SuperTrend`, `UlcerIndex`.

## Pattern

Volatility indicators typically return channels of complex types (e.g., `BollingerBands` returns `<-chan Bands[T]`) or simple types like `Atr`.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `GEMINI.md`), using `testdata/*.csv` files (e.g., `atr.csv`, `bollinger_bands.csv`) for historical data validation.
