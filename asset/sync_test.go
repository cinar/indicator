// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

type MockRepository struct {
	asset.Repository

	AssetsFunc   func() ([]string, error)
	GetFunc      func(string) (<-chan *asset.Snapshot, error)
	GetSinceFunc func(string, time.Time) (<-chan *asset.Snapshot, error)
	LastDateFunc func(string) (time.Time, error)
	AppendFunc   func(string, <-chan *asset.Snapshot) error
}

func (r *MockRepository) Assets() ([]string, error) {
	return r.AssetsFunc()
}

func (r *MockRepository) Get(name string) (<-chan *asset.Snapshot, error) {
	return r.GetFunc(name)
}

func (r *MockRepository) GetSince(name string, date time.Time) (<-chan *asset.Snapshot, error) {
	return r.GetSinceFunc(name, date)
}

func (r *MockRepository) LastDate(name string) (time.Time, error) {
	return r.LastDateFunc(name)
}

func (r *MockRepository) Append(name string, snapshots <-chan *asset.Snapshot) error {
	return r.AppendFunc(name, snapshots)
}

func TestSync(t *testing.T) {
	name := "A"
	snapshots := []*asset.Snapshot{
		{Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	source := asset.NewInMemoryRepository()
	target := asset.NewInMemoryRepository()

	err := target.Append(name, helper.SliceToChan([]*asset.Snapshot{}))
	if err != nil {
		t.Fatal(err)
	}

	err = source.Append(name, helper.SliceToChan(snapshots))
	if err != nil {
		t.Fatal(err)
	}

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err = sync.Run(source, target, snapshots[0].Date)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := target.Get(name)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSyncMissingOnSource(t *testing.T) {
	name := "A"

	source := asset.NewInMemoryRepository()
	target := asset.NewInMemoryRepository()

	err := target.Append(name, helper.SliceToChan([]*asset.Snapshot{}))
	if err != nil {
		t.Fatal(err)
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err = sync.Run(source, target, defaultStartDate)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSyncFailingTargetAssets(t *testing.T) {
	source := asset.NewInMemoryRepository()
	target := &MockRepository{
		AssetsFunc: func() ([]string, error) {
			return nil, errors.New("assert error")
		},
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err := sync.Run(source, target, defaultStartDate)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSyncFailingTargetAppendMultiWorker(t *testing.T) {
	assetNames := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	source := &MockRepository{
		GetSinceFunc: func(_ string, _ time.Time) (<-chan *asset.Snapshot, error) {
			return helper.SliceToChan([]*asset.Snapshot{}), nil
		},
	}

	target := &MockRepository{
		AssetsFunc: func() ([]string, error) {
			return assetNames, nil
		},

		LastDateFunc: func(_ string) (time.Time, error) {
			return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil
		},

		AppendFunc: func(_ string, _ <-chan *asset.Snapshot) error {
			return errors.New("append error")
		},
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 4
	sync.Delay = 0

	err := sync.Run(source, target, defaultStartDate)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSyncLastDateRealErrorDoesNotFallBackToDefault(t *testing.T) {
	name := "A"

	var getSinceCalledWithDate time.Time
	getSinceCalled := false

	source := &MockRepository{
		GetSinceFunc: func(_ string, date time.Time) (<-chan *asset.Snapshot, error) {
			getSinceCalled = true
			getSinceCalledWithDate = date

			return helper.SliceToChan([]*asset.Snapshot{}), nil
		},
	}

	target := &MockRepository{
		AssetsFunc: func() ([]string, error) {
			return []string{name}, nil
		},

		LastDateFunc: func(_ string) (time.Time, error) {
			return time.Time{}, errors.New("db connection dropped")
		},

		AppendFunc: func(_ string, snapshots <-chan *asset.Snapshot) error {
			helper.Drain(snapshots)
			return nil
		},
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err := sync.Run(source, target, defaultStartDate)
	if err == nil {
		t.Fatal("expected error")
	}

	if getSinceCalled {
		t.Fatalf("expected GetSince to not be called, but it was called with date %v", getSinceCalledWithDate)
	}
}

func TestSyncLastDateSentinelErrorFallsBackToDefault(t *testing.T) {
	name := "A"

	var getSinceCalledWithDate time.Time

	source := &MockRepository{
		GetSinceFunc: func(_ string, date time.Time) (<-chan *asset.Snapshot, error) {
			getSinceCalledWithDate = date

			return helper.SliceToChan([]*asset.Snapshot{}), nil
		},
	}

	target := &MockRepository{
		AssetsFunc: func() ([]string, error) {
			return []string{name}, nil
		},

		LastDateFunc: func(_ string) (time.Time, error) {
			return time.Time{}, asset.ErrRepositoryAssetEmpty
		},

		AppendFunc: func(_ string, snapshots <-chan *asset.Snapshot) error {
			helper.Drain(snapshots)
			return nil
		},
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err := sync.Run(source, target, defaultStartDate)
	if err != nil {
		t.Fatal(err)
	}

	if !getSinceCalledWithDate.Equal(defaultStartDate) {
		t.Fatalf("expected GetSince to be called with default start date %v, got %v", defaultStartDate, getSinceCalledWithDate)
	}
}

func TestSyncFailingTargetAppend(t *testing.T) {
	source := &MockRepository{
		GetSinceFunc: func(_ string, _ time.Time) (<-chan *asset.Snapshot, error) {
			return helper.SliceToChan([]*asset.Snapshot{}), nil
		},
	}

	target := &MockRepository{
		AssetsFunc: func() ([]string, error) {
			return []string{"A"}, nil
		},

		LastDateFunc: func(_ string) (time.Time, error) {
			return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil
		},

		AppendFunc: func(_ string, _ <-chan *asset.Snapshot) error {
			return errors.New("append error")
		},
	}

	defaultStartDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	sync := asset.NewSync()
	sync.Workers = 1
	sync.Delay = 0

	err := sync.Run(source, target, defaultStartDate)
	if err == nil {
		t.Fatal("expected error")
	}
}
