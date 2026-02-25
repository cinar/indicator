# Helper Package - Utilities

The `helper` package contains essential low-level utilities for channel manipulation, math operations, and testing support used throughout the project.

## Key Utilities

- **Channels:** `SliceToChan`, `ChanToSlice`, `Buffered`, `Duplicate`, `Head`, `Last`, `Skip`, `Pipe`.
- **Math:** `Abs`, `Add`, `Divide`, `Multiply`, `Subtract`, `Pow`, `Sqrt`, `Sign`, `RoundDigit`, `RoundDigits`.
- **Stats:** `Highest`, `Lowest`, `Since`, `MaxSince`, `MinSince`, `DaysBetween`.
- **Data:** `CsvToChan`, `ChanToCsv`, `JsonToChan`, `ChanToJson`, `ReadFromCsvFile`.

## Number Interface

All helpers use the `helper.Number` interface which encompasses `int`, `float32`, and `float64`.

```go
type Number interface {
	Integer | Float
}
```

## Testing Helper

- `CheckEquals`: Basic comparison for test values.
- `RoundDigits`: Essential for comparing floating-point results.
- `ReadFromCsvFile`: Load large test datasets from `testdata/`.

## Pattern

Helpers are often implemented as standalone functions or channel transformations:
```go
func SliceToChan[T any](s []T) <-chan T { ... }
func ChanToSlice[T any](c <-chan T) []T { ... }
```
