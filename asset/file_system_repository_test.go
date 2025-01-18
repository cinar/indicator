// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"fmt"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

var repositoryBase = "testdata/repository"

func TestFileSystemRepositoryAssets(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)
	expected := []string{"brk-b"}

	actual, err := repository.Assets()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestFileSystemRepositoryAssetsNonExisting(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/non_existing")

	_, err := repository.Assets()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFileSystemRepositoryGet(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	snapshots, err := repository.Get("brk-b")
	if err != nil {
		t.Fatal(err)
	}

	helper.Drain(snapshots)
}

func TestFileSystemRepositoryGetNonExisting(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/non_existing")

	_, err := repository.Get("brk-b")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFileSystemRepositoryGetSince(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	date := time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
	actual, err := repository.GetSince("brk-b", date)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/since.csv")
	if err != nil {
		t.Fatal(err)
	}

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileSystemRepositoryGetSinceNonExisting(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	date := time.Date(2022, 12, 0o1, 0, 0, 0, 0, time.UTC)
	_, err := repository.GetSince("brk", date)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFileSystemRepositoryLastDate(t *testing.T) {
	expeted := time.Date(2023, 11, 29, 0, 0, 0, 0, time.UTC)

	repository := asset.NewFileSystemRepository(repositoryBase)

	actual, err := repository.LastDate("brk-b")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expeted {
		t.Fatalf("actual %v expected %v", actual, expeted)
	}
}

func TestFileSystemRepositoryLastDateNonExisting(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	_, err := repository.LastDate("brk")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFileSystemRepositoryLastDateEmpty(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	_, err := repository.LastDate("empty")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFileSystemRepositoryAppend(t *testing.T) {
	repository := asset.NewFileSystemRepository(repositoryBase)

	expected, err := repository.Get("brk-b")
	if err != nil {
		t.Fatal(err)
	}

	name := "test_file_system_repository_append"
	defer helper.Remove(t, path.Join(repositoryBase, fmt.Sprintf("%s.csv", name)))

	err = repository.Append(name, expected)
	if err != nil {
		t.Fatal(err)
	}

	expected, err = repository.Get("brk-b")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := repository.Get(name)
	if err != nil {
		t.Fatal(err)
	}

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
