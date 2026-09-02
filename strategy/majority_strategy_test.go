// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestMajorityStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/majority.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	majority := strategy.NewMajorityStrategy("Majority Strategy")
	majority.Strategies = append(majority.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	actual := majority.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMajorityStrategyNoStrategies(t *testing.T) {
	snapshots := helper.SliceToChan([]*asset.Snapshot{})

	majority := strategy.NewMajorityStrategy("Majority Strategy")
	actual := majority.Compute(snapshots)

	select {
	case _, ok := <-actual:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout - result channel never closed")
	}
}

func TestMajorityStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	majority := strategy.NewMajorityStrategy("Majority Strategy")
	majority.Strategies = append(majority.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	report := majority.Report(snapshots)

	fileName := "majority.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
