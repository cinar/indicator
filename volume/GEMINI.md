# Volume Package - Indicators

The `volume` package provides a collection of indicators for analyzing price-volume relationships.

## Key Indicators

- **Price-Volume:** `Adi` (Accumulation Distribution Index), `Adx` (Average Directional Index), `Cmfi` (Chaikin Money Flow Index), `Emv` (Ease of Movement).
- **Oscillators:** `Kvo` (Klinger Volume Oscillator), `Obv` (On Balance Volume), `Pvi` (Positive Volume Index), `Nvi` (Negative Volume Index).
- **Indicators:** `Mfi` (Money Flow Index), `VortexIndicator`.

## Pattern

Volume indicators use both price and volume data streams, often via `helper.Number` generics and multiple channel inputs.

## Testing Standard

Tests in this package follow the project's standard CSV-based testing pattern (as detailed in the root `GEMINI.md`), using `testdata/*.csv` files (e.g., `obv.csv`, `mfi.csv`) for historical data validation.
