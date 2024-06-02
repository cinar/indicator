// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"os"
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestBacktest(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	outputDir, err := os.MkdirTemp("", "backtest")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(outputDir)

	backtest := strategy.NewBacktest(repository, outputDir)
	backtest.Names = append(backtest.Names, "brk-b")
	backtest.Strategies = append(backtest.Strategies, trend.NewApoStrategy())

	err = backtest.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBacktestAllAssetsAndStrategies(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	outputDir, err := os.MkdirTemp("", "backtest")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(outputDir)

	backtest := strategy.NewBacktest(repository, outputDir)

	err = backtest.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBacktestNonExistingAsset(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	outputDir, err := os.MkdirTemp("", "backtest")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(outputDir)

	backtest := strategy.NewBacktest(repository, outputDir)
	backtest.Names = append(backtest.Names, "non_existing")

	err = backtest.Run()
	if err != nil {
		t.Fatal(err)
	}
}
