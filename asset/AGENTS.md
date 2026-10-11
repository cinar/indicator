# Asset Package - Data Management

The `asset` package handles asset data management and persistence, supporting various repository types for historical data.

## Key Components

- **Models:** `Snapshot` (a single OHLCV bar). Assets are identified by plain `string` names, not a separate `Asset` struct.
- **Repositories:** `Repository`, `FileSystemRepository`, `InMemoryRepository`, `SQLRepository`, `TiingoRepository`.
- **Factory:** `NewRepository(name, config string) (Repository, error)`, backed by `RegisterRepositoryBuilder` for adding new repository types.
- **Utilities:** `Sync` (keeps a repository's snapshots up to date from a remote source).

## Repository Pattern

The `Repository` interface defines common operations for interacting with asset data.
```go
type Repository interface {
	Assets() ([]string, error)
	Get(name string) (<-chan *Snapshot, error)
	GetSince(name string, date time.Time) (<-chan *Snapshot, error)
	LastDate(name string) (time.Time, error)
	Append(name string, snapshots <-chan *Snapshot) error
}
```

## Storage Types

- `InMemoryRepository`: Fast for ephemeral data and testing.
- `FileSystemRepository`: Storage on disk using CSV or JSON files.
- `SQLRepository`: Database-backed persistence for large datasets.
- `TiingoRepository`: Remote API connector for fetching historical data from Tiingo.

## Model Consistency

All price and volume fields in repository models (e.g., `TiingoEndOfDay`) must use `float64`. This ensures compatibility with crypto assets that provide fractional volumes, which would otherwise cause JSON unmarshaling errors if `int64` is used.

## Testing Pattern

Test files use the `asset_test` package and verify repository implementations against mock and real-world data sources.
```go
func TestInMemoryRepository(t *testing.T) {
	repo := asset.NewInMemoryRepository()
	// test Append, Get, GetSince, LastDate, etc.
}
```
