# Trend Package - Indicators

The `trend` package provides a collection of indicators used for identifying and analyzing trends in price data.

## Key Indicators

- **Moving Averages:** `Sma` (Simple), `Ema` (Exponential), `Dema` (Double), `Tema` (Triple), `Wma` (Weighted), `Hma` (Hull), `Kama` (Kaufman), `Smma` (Smoothed), `Vwma` (Volume Weighted).
- **Oscillators:** `Apo` (Absolute Price Oscillator), `Cci` (Commodity Channel Index), `Dpo` (Detrended Price Oscillator), `Trix` (Triple Exponential Average).
- **Indicators:** `Aroon` (Aroon Oscillator), `Bop` (Balance of Power), `Macd` (Moving Average Convergence Divergence), `Roc` (Rate of Change), `Tsi` (True Strength Index).
- **Utilities:** `MovingMax`, `MovingMin`, `MovingSum`, `TypicalPrice`.

## Common Pattern

All indicators follow the standard `Compute(<-chan T) <-chan T` pattern and are generic over `helper.Number`.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `GEMINI.md`), using `testdata/*.csv` files (e.g., `macd.csv`, `ema.csv`) for historical data validation.
