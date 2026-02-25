# Valuation Package - Indicators

The `valuation` package provides indicators for price valuation and statistical analysis.

## Key Indicators

- **Price-Earnings:** `PeRatio` (Price to Earnings Ratio), `PsRatio` (Price to Sales Ratio), `PbRatio` (Price to Book Ratio).
- **Yields:** `DividendYield`, `EarningsYield`.
- **Metrics:** `EnterpriseValue`, `MarketCap`.

## Pattern

Valuation indicators often involve multi-stream data like price and financial metrics, processed via `Compute(<-chan T, <-chan T) <-chan T`.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `GEMINI.md`), using `testdata/*.csv` files for historical data validation.
