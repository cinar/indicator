// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
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

func TestBacktestNoStrategies(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	outputDir, err := os.MkdirTemp("", "backtest")
	if err != nil {
		t.Fatal(err)
	}

	defer helper.RemoveAll(t, outputDir)

	htmlReport := backtest.NewHTMLReport(outputDir)
	bt := backtest.NewBacktest(repository, htmlReport)
	bt.Names = append(bt.Names, "brk-b")

	err = bt.Run()
	if err != nil {
		t.Fatal(err)
	}
}

type errorReport struct {
	backtest.Report
}

func (e *errorReport) Begin([]string, []strategy.Strategy) error {
	return fmt.Errorf("begin error")
}

func TestBacktestBeginError(t *testing.T) {
	repository := asset.NewInMemoryRepository()
	report := &errorReport{}
	bt := backtest.NewBacktest(repository, report)

	err := bt.Run()
	if err == nil {
		t.Fatal("expected error")
	}
}

type assetErrorReport struct {
	backtest.Report
}

func (e *assetErrorReport) Begin([]string, []strategy.Strategy) error { return nil }
func (e *assetErrorReport) AssetBegin(string, []strategy.Strategy) error {
	return fmt.Errorf("asset begin error")
}
func (e *assetErrorReport) End() error { return fmt.Errorf("end error") }

func TestBacktestAssetBeginError(t *testing.T) {
	repository := asset.NewInMemoryRepository()
	snapshots := make(chan *asset.Snapshot, 1)
	snapshots <- &asset.Snapshot{Date: time.Now()}
	close(snapshots)
	repository.Append("TEST", snapshots)

	report := &assetErrorReport{}
	bt := backtest.NewBacktest(repository, report)
	bt.Names = []string{"TEST"}

	err := bt.Run()
	if err == nil {
		t.Fatal("expected error from End()")
	}
}
