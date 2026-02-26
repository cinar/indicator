// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestHTMLReport(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "report")
	if err != nil {
		t.Fatal(err)
	}
	defer helper.RemoveAll(t, outputDir)

	report := backtest.NewHTMLReport(outputDir)

	assetNames := []string{"TEST"}
	strategies := []strategy.Strategy{strategy.NewBuyAndHoldStrategy()}

	err = report.Begin(assetNames, strategies)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetBegin("TEST", strategies)
	if err != nil {
		t.Fatal(err)
	}

	snapshots := make(chan *asset.Snapshot, 1)
	snapshots <- &asset.Snapshot{Close: 100}
	close(snapshots)

	actions := make(chan strategy.Action, 1)
	actions <- strategy.Buy
	close(actions)

	outcomes := make(chan float64, 1)
	outcomes <- 1.1
	close(outcomes)

	err = report.Write("TEST", strategies[0], snapshots, actions, outcomes)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetEnd("TEST")
	if err != nil {
		t.Fatal(err)
	}

	err = report.End()
	if err != nil {
		t.Fatal(err)
	}

	// Check if files exist
	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); os.IsNotExist(err) {
		t.Fatal("index.html not found")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "TEST.html")); os.IsNotExist(err) {
		t.Fatal("TEST.html not found")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "TEST - Buy and Hold Strategy.html")); os.IsNotExist(err) {
		t.Fatal("strategy report not found")
	}
}

func TestHTMLReportErrors(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "report_err")
	if err != nil {
		t.Fatal(err)
	}
	defer helper.RemoveAll(t, outputDir)

	report := backtest.NewHTMLReport(outputDir)

	err = report.AssetBegin("TEST", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetBegin("TEST", nil)
	if err == nil {
		t.Fatal("expected error for already begun asset")
	}

	snapshots := make(chan *asset.Snapshot, 1)
	snapshots <- &asset.Snapshot{Close: 100}
	close(snapshots)

	actions := make(chan strategy.Action, 1)
	actions <- strategy.Buy
	close(actions)

	outcomes := make(chan float64, 1)
	outcomes <- 1.1
	close(outcomes)

	actions = make(chan strategy.Action, 1)
	actions <- strategy.Buy
	close(actions)

	outcomes = make(chan float64, 1)
	outcomes <- 1.1
	close(outcomes)

	go helper.Drain(snapshots)

	err = report.Write("UNKNOWN", strategy.NewBuyAndHoldStrategy(), snapshots, actions, outcomes)
	if err == nil {
		t.Fatal("expected error for not begun asset")
	}

	err = report.AssetEnd("UNKNOWN")
	if err == nil {
		t.Fatal("expected error for not begun asset")
	}
}

func TestHTMLReportNoStrategyReports(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "report_no_strat")
	if err != nil {
		t.Fatal(err)
	}
	defer helper.RemoveAll(t, outputDir)

	report := backtest.NewHTMLReport(outputDir)
	report.WriteStrategyReports = false

	assetNames := []string{"TEST"}
	strategies := []strategy.Strategy{strategy.NewBuyAndHoldStrategy()}

	err = report.Begin(assetNames, strategies)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetBegin("TEST", strategies)
	if err != nil {
		t.Fatal(err)
	}

	snapshots := make(chan *asset.Snapshot, 1)
	snapshots <- &asset.Snapshot{Close: 100}
	close(snapshots)

	actions := make(chan strategy.Action, 1)
	actions <- strategy.Buy
	close(actions)

	outcomes := make(chan float64, 1)
	outcomes <- 1.1
	close(outcomes)

	err = report.Write("TEST", strategies[0], snapshots, actions, outcomes)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetEnd("TEST")
	if err != nil {
		t.Fatal(err)
	}

	err = report.End()
	if err != nil {
		t.Fatal(err)
	}
}
