// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

func TestInMemoryRepositoryAssets(t *testing.T) {
	repository := asset.NewInMemoryRepository()

	assets, err := repository.Assets()
	if err != nil {
		t.Fatal(err)
	}

	if len(assets) != 0 {
		t.Fatal("not empty")
	}

	name := "A"

	snapshots := []*asset.Snapshot{
		{Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	err = repository.Append(name, helper.SliceToChan(snapshots))
	if err != nil {
		t.Fatal(err)
	}

	assets, err = repository.Assets()
	if err != nil {
		t.Fatal(err)
	}

	if len(assets) != 1 {
		t.Fatalf("more assets found %v", assets)
	}

	if assets[0] != name {
		t.Fatalf("actual %v expected %v", assets[0], name)
	}
}

func TestInMemoryRepositoryGet(t *testing.T) {
	repository := asset.NewInMemoryRepository()

	name := "A"

	_, err := repository.Get(name)
	if err == nil {
		t.Fatal("expected error")
	}

	snapshots := []*asset.Snapshot{
		{Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	err = repository.Append(name, helper.SliceToChan(snapshots))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := repository.Get(name)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInMemoryRepositoryGetSince(t *testing.T) {
	repository := asset.NewInMemoryRepository()

	name := "A"
	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)

	_, err := repository.GetSince(name, date)
	if err == nil {
		t.Fatal("expected error")
	}

	snapshots := []*asset.Snapshot{
		{Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	err = repository.Append(name, helper.SliceToChan(snapshots))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := repository.GetSince(name, date)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(snapshots[1:])

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInMemoryRepositoryLastDate(t *testing.T) {
	repository := asset.NewInMemoryRepository()

	name := "A"

	_, err := repository.LastDate(name)
	if err == nil {
		t.Fatal("expected error")
	}

	err = repository.Append(name, helper.SliceToChan([]*asset.Snapshot{}))
	if err != nil {
		t.Fatal(err)
	}

	_, err = repository.LastDate(name)
	if err == nil {
		t.Fatal("expected error")
	}

	snapshots := []*asset.Snapshot{
		{Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	err = repository.Append(name, helper.SliceToChan(snapshots))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := repository.LastDate(name)
	if err != nil {
		t.Fatal(err)
	}

	expected := snapshots[1].Date

	if !expected.Equal(actual) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestInMemoryRepositoryConcurrentAccess(t *testing.T) {
	repository := asset.NewInMemoryRepository()

	const (
		workers        = 8
		appendsPerName = 50
	)

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		name := fmt.Sprintf("ASSET_%d", w)

		wg.Add(2)

		// Concurrently appends snapshots to the same asset name.
		go func(name string) {
			defer wg.Done()

			for i := 0; i < appendsPerName; i++ {
				snapshot := &asset.Snapshot{
					Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i),
				}

				err := repository.Append(name, helper.SliceToChan([]*asset.Snapshot{snapshot}))
				if err != nil {
					t.Error(err)
				}
			}
		}(name)

		// Concurrently reads the storage while the appends above are in flight.
		go func(name string) {
			defer wg.Done()

			for i := 0; i < appendsPerName; i++ {
				if _, err := repository.Assets(); err != nil {
					t.Error(err)
				}

				channel, err := repository.Get(name)
				if err != nil {
					continue
				}

				for range channel {
				}
			}
		}(name)
	}

	wg.Wait()
}
