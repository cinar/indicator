// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"errors"
	"time"
)

// ErrRepositoryAssetNotFound indicates that the given asset name is not found in the repository.
var ErrRepositoryAssetNotFound = errors.New("asset is not found")

// ErrRepositoryAssetEmpty indicates that the given asset has no snapshots.
var ErrRepositoryAssetEmpty = errors.New("asset empty")

// Repository serves as a centralized storage and retrieval
// location for asset snapshots.
type Repository interface {
	// Assets returns the names of all assets in the repository.
	Assets() ([]string, error)

	// Get attempts to return a channel of snapshots for
	// the asset with the given name.
	//
	// Some implementations (e.g. SQLRepository, TiingoRepository) only
	// return snapshots from 2000-01-01 onward by design, while others
	// (e.g. InMemoryRepository, FileSystemRepository) return full
	// history with no floor date; see each implementation's doc comment.
	Get(name string) (<-chan *Snapshot, error)

	// GetSince attempts to return a channel of snapshots for
	// the asset with the given name since the given date.
	//
	// Implementations may run the underlying retrieval in a background
	// goroutine that sends on the returned channel; there is no
	// cancellation parameter on this method, so callers must read the
	// returned channel to completion. Abandoning it partway through can
	// leak the goroutine and any resources it holds (e.g. a database
	// connection or an HTTP response body) for implementations that rely
	// on the channel being drained to release them.
	GetSince(name string, date time.Time) (<-chan *Snapshot, error)

	// LastDate returns the date of the last snapshot for
	// the asset with the given name.
	LastDate(name string) (time.Time, error)

	// Append adds the given snapshows to the asset with the
	// given name.
	Append(name string, snapshots <-chan *Snapshot) error
}
