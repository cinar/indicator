// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"github.com/cinar/indicator/v2/helper"
	"os"
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestBacktest(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	outputDir, err := os.MkdirTemp("", "bt")
	if err != nil {
		t.Fatal(err)
	}

	defer helper.RemoveAll(t, outputDir)

	htmlReport := backtest.NewHTMLReport(outputDir)
	bt := backtest.NewBacktest(repository, htmlReport)
	bt.Names = append(bt.Names, "brk-b")
	bt.Strategies = append(bt.Strategies, trend.NewApoStrategy())

	err = bt.Run()
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

	defer helper.RemoveAll(t, outputDir)

	htmlReport := backtest.NewHTMLReport(outputDir)
	bt := backtest.NewBacktest(repository, htmlReport)

	err = bt.Run()
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

	defer helper.RemoveAll(t, outputDir)

	htmlReport := backtest.NewHTMLReport(outputDir)
	bt := backtest.NewBacktest(repository, htmlReport)
	bt.Names = append(bt.Names, "non_existing")

	err = bt.Run()
	if err != nil {
		t.Fatal(err)
	}
}
