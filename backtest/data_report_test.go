// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestDataReport(t *testing.T) {
	repository := asset.NewFileSystemRepository("testdata/repository")

	assets := []string{
		"brk-b",
	}

	strategies := []strategy.Strategy{
		strategy.NewBuyAndHoldStrategy(),
	}

	report := backtest.NewDataReport()

	err := report.Begin(assets, strategies)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetBegin(assets[0], strategies)
	if err != nil {
		t.Fatal(err)
	}

	snapshots, err := repository.Get(assets[0])
	if err != nil {
		t.Fatal(err)
	}

	snapshotsSplice := helper.Duplicate(snapshots, 3)
	actionsSplice := helper.Duplicate(
		strategies[0].Compute(snapshotsSplice[1]),
		2,
	)

	outcomes := strategy.Outcome(
		asset.SnapshotsAsClosings(snapshotsSplice[2]),
		actionsSplice[1],
	)

	err = report.Write(assets[0], strategies[0], snapshotsSplice[0], actionsSplice[0], outcomes)
	if err != nil {
		t.Fatal(err)
	}

	err = report.AssetEnd(assets[0])
	if err != nil {
		t.Fatal(err)
	}

	err = report.End()
	if err != nil {
		t.Fatal(err)
	}

	results, ok := report.Results[assets[0]]
	if !ok {
		t.Fatal("asset result not found")
	}

	if len(results) != len(strategies) {
		t.Fatalf("results count and strategies count are not the same, %d %d", len(results), len(strategies))
	}
}
