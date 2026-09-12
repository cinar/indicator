# Volume Package - Indicators

The `volume` package provides a collection of indicators for analyzing price-volume relationships.

## Key Indicators

- **Price-Volume:** `Ad` (Accumulation/Distribution), `Cmf` (Chaikin Money Flow), `Mfm`/`Mfv` (Money Flow Multiplier/Volume), `Emv` (Ease of Movement), `Vwap` (Volume Weighted Average Price).
- **Oscillators:** `Kvo` (Klinger Volume Oscillator), `Obv` (On Balance Volume), `Nvi` (Negative Volume Index).
- **Indicators:** `Mfi` (Money Flow Index), `Fi` (Force Index), `Vpt` (Volume Price Trend).

## Pattern

Volume indicators use both price and volume data streams, generic over `helper.Float` (`Nvi`, `Obv`, `Vpt` use the broader `helper.Number`), with multiple channel inputs.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `AGENTS.md`), using `testdata/*.csv` files (e.g., `obv.csv`, `mfi.csv`) for historical data validation.
