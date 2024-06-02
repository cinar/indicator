// Copyright (c) 2021-2024 Onur Cinar.
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
	Get(name string) (<-chan *Snapshot, error)

	// GetSince attempts to return a channel of snapshots for
	// the asset with the given name since the given date.
	GetSince(name string, date time.Time) (<-chan *Snapshot, error)

	// LastDate returns the date of the last snapshot for
	// the asset with the given name.
	LastDate(name string) (time.Time, error)

	// Append adds the given snapshows to the asset with the
	// given name.
	Append(name string, snapshots <-chan *Snapshot) error
}
