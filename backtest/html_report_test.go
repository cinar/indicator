// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"os"
	"path/filepath"
	"strings"
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

func TestHTMLReportSortsCloseOutcomes(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "report_sort")
	if err != nil {
		t.Fatal(err)
	}
	defer helper.RemoveAll(t, outputDir)

	report := backtest.NewHTMLReport(outputDir)
	report.WriteStrategyReports = false

	strategies := []strategy.Strategy{
		strategy.NewMajorityStrategy("Strategy A"),
		strategy.NewMajorityStrategy("Strategy B"),
	}

	err = report.Begin([]string{"TEST"}, strategies)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetBegin("TEST", strategies)
	if err != nil {
		t.Fatal(err)
	}

	// Outcomes 1.2 and 1.6 (0.012 and 0.016 before the *100 scaling in
	// Write) differ by less than 1.0, so int(b-a) truncation would treat
	// them as equal. They must still sort with the 1.6 outcome first.
	writeOutcome := func(currentStrategy strategy.Strategy, outcome float64) {
		snapshots := make(chan *asset.Snapshot, 1)
		snapshots <- &asset.Snapshot{Close: 100}
		close(snapshots)

		actions := make(chan strategy.Action, 1)
		actions <- strategy.Buy
		close(actions)

		outcomes := make(chan float64, 1)
		outcomes <- outcome
		close(outcomes)

		err = report.Write("TEST", currentStrategy, snapshots, actions, outcomes)
		if err != nil {
			t.Fatal(err)
		}
	}

	writeOutcome(strategies[0], 0.012)
	writeOutcome(strategies[1], 0.016)

	err = report.AssetEnd("TEST")
	if err != nil {
		t.Fatal(err)
	}

	err = report.End()
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "TEST.html"))
	if err != nil {
		t.Fatal(err)
	}

	indexA := strings.Index(string(content), "Strategy A")
	indexB := strings.Index(string(content), "Strategy B")
	if indexA == -1 || indexB == -1 {
		t.Fatalf("expected both strategies in report, got: %s", content)
	}

	if indexB > indexA {
		t.Fatalf("expected higher outcome (Strategy B, 1.6) to be listed before lower outcome (Strategy A, 1.2)")
	}
}

func TestHTMLReportAssetEndNoResults(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "report_no_results")
	if err != nil {
		t.Fatal(err)
	}
	defer helper.RemoveAll(t, outputDir)

	report := backtest.NewHTMLReport(outputDir)

	err = report.AssetBegin("TEST", nil)
	if err != nil {
		t.Fatal(err)
	}

	// No strategies were written for the asset, so AssetEnd must not panic
	// on an empty results slice.
	err = report.AssetEnd("TEST")
	if err != nil {
		t.Fatalf("expected no error for asset with no results, got: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "TEST.html")); os.IsNotExist(err) {
		t.Fatal("TEST.html not found")
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
