// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
)

func TestNewRepositoryUnknown(t *testing.T) {
	repository, err := asset.NewRepository("unknown", "")
	if err == nil {
		t.Fatalf("unknown repository: %T", repository)
	}
}

func TestRegisterRepositoryBuilder(t *testing.T) {
	builderName := "testbuilder"

	repository, err := asset.NewRepository(builderName, "")
	if err == nil {
		t.Fatalf("testbuilder is: %T", repository)
	}

	asset.RegisterRepositoryBuilder(builderName, func(_ string) (asset.Repository, error) {
		return asset.NewInMemoryRepository(), nil
	})

	repository, err = asset.NewRepository(builderName, "")
	if err != nil {
		t.Fatalf("testbuilder is not found: %v", err)
	}

	_, ok := repository.(*asset.InMemoryRepository)
	if !ok {
		t.Fatalf("testbuilder is: %T", repository)
	}
}

func TestNewRepositoryMemory(t *testing.T) {
	repository, err := asset.NewRepository(asset.InMemoryRepositoryBuilderName, "")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := repository.(*asset.InMemoryRepository)
	if !ok {
		t.Fatalf("repository not correct type: %T", repository)
	}
}

func TestNewRepositoryFileSystem(t *testing.T) {
	repository, err := asset.NewRepository(asset.FileSystemRepositoryBuilderName, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := repository.(*asset.FileSystemRepository)
	if !ok {
		t.Fatalf("repository not correct type: %T", repository)
	}
}

func TestNewTiingoRepository(t *testing.T) {
	repository, err := asset.NewRepository(asset.TiingoRepositoryBuilderName, "1234")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := repository.(*asset.TiingoRepository)
	if !ok {
		t.Fatalf("repository not correct type: %T", repository)
	}
}
