// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"sync"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

// InMemoryRepository stores and retrieves asset snapshots using
// an in memory storage. It is safe for concurrent use.
type InMemoryRepository struct {
	// mutex guards access to storage.
	mutex sync.RWMutex

	// storage is the in memory storage for assets.
	storage map[string][]*Snapshot
}

// NewInMemoryRepository initializes an in memory repository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(map[string][]*Snapshot),
	}
}

// Assets returns the names of all assets in the repository.
func (r *InMemoryRepository) Assets() ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	assets := make([]string, 0, len(r.storage))
	for name := range r.storage {
		assets = append(assets, name)
	}

	return assets, nil
}

// Get attempts to return a channel of snapshots for the asset with the given name.
func (r *InMemoryRepository) Get(name string) (<-chan *Snapshot, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	snapshots, ok := r.storage[name]
	if !ok {
		return nil, ErrRepositoryAssetNotFound
	}

	// Copy while holding the lock, since SliceToChan reads the slice
	// lazily from a goroutine, after a concurrent Append may have
	// mutated this slice's backing array in place.
	copied := make([]*Snapshot, len(snapshots))
	copy(copied, snapshots)

	return helper.SliceToChan(copied), nil
}

// GetSince attempts to return a channel of snapshots for the asset with the given name since the given date.
func (r *InMemoryRepository) GetSince(name string, date time.Time) (<-chan *Snapshot, error) {
	snapshots, err := r.Get(name)
	if err != nil {
		return nil, err
	}

	snapshots = helper.Filter(snapshots, func(s *Snapshot) bool {
		return s.Date.Equal(date) || s.Date.After(date)
	})

	return snapshots, nil
}

// LastDate returns the date of the last snapshot for the asset with the given name.
func (r *InMemoryRepository) LastDate(name string) (time.Time, error) {
	var last time.Time

	snapshots, err := r.Get(name)
	if err != nil {
		return last, err
	}

	snapshot, ok := <-helper.Last(snapshots, 1)
	if !ok {
		return last, ErrRepositoryAssetEmpty
	}

	return snapshot.Date, nil
}

// Append adds the given snapshows to the asset with the given name.
func (r *InMemoryRepository) Append(name string, snapshots <-chan *Snapshot) error {
	var newSnapshots []*Snapshot

	for snapshot := range snapshots {
		newSnapshots = append(newSnapshots, snapshot)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.storage[name] = append(r.storage[name], newSnapshots...)

	return nil
}
