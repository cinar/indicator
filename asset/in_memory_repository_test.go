// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
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
