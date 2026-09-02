// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/examples/trend"
)

func TestMacdStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/macd_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	macd := trend.NewMacdStrategy()
	actual := macd.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMacdStrategyEdgeTriggered(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/macd_strategy_edge_triggered.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.ChanToSlice(helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action }))

	macd := trend.NewMacdStrategy()
	macd.SignalMode = trend.EdgeTriggered

	actual := macd.Compute(snapshots)

	err = helper.CheckEquals(actual, helper.SliceToChan(expected))
	if err != nil {
		t.Fatal(err)
	}

	// EdgeTriggered fires only once, at the crossing bar, so it should produce far
	// fewer non-Hold actions than the default LevelTriggered mode (which fires on
	// every bar that still satisfies the crossover condition).
	buys, sells := 0, 0
	for _, action := range expected {
		switch action {
		case strategy.Buy:
			buys++
		case strategy.Sell:
			sells++
		}
	}

	const wantBuys, wantSells = 4, 8

	if buys != wantBuys || sells != wantSells {
		t.Fatalf("got %d buys and %d sells, want %d buys and %d sells", buys, sells, wantBuys, wantSells)
	}
}

func TestMacdStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	macd := trend.NewMacdStrategy()

	report := macd.Report(snapshots)

	fileName := "macd_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
