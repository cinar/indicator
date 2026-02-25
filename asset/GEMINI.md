# Asset Package - Data Management

The `asset` package handles asset data management and persistence, supporting various repository types for historical data.

## Key Components

- **Models:** `Asset`, `Snapshot`, `SnapshotType`.
- **Repositories:** `Repository`, `FileSystemRepository`, `InMemoryRepository`, `SqlRepository`, `TiingoRepository`.
- **Factories:** `RepositoryFactory`, `RepositoryConfig`.
- **Utilities:** `Sync`.

## Repository Pattern

The `Repository` interface defines common operations for interacting with asset data.
```go
type Repository interface {
	GetAssets() ([]Asset, error)
	GetSnapshots(asset string) (<-chan Snapshot, error)
	AddSnapshot(asset string, snapshot Snapshot) error
}
```

## Storage Types

- `InMemoryRepository`: Fast for ephemeral data and testing.
- `FileSystemRepository`: Storage on disk using CSV or JSON files.
- `SqlRepository`: Database-backed persistence for large datasets.
- `TiingoRepository`: Remote API connector for fetching real-time data.

## Testing Pattern

Test files use `asset_test` package and verify repository implementations against mock and real-world data sources.
```go
func TestInMemoryRepository(t *testing.T) {
	repo := asset.NewInMemoryRepository()
	// test AddSnapshot, GetSnapshots, etc.
}
```
