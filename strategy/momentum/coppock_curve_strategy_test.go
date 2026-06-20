// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"context"
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	indMomentum "github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/momentum"
)

func TestCoppockCurveStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/coppock_curve_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	cc := momentum.NewCoppockCurveStrategy()
	actual := cc.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCoppockCurveStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	cc := momentum.NewCoppockCurveStrategy()

	report := cc.Report(snapshots)

	fileName := "coppock_curve_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCoppockCurveStrategyZeroAndContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := momentum.NewCoppockCurveStrategy()
	cc.CoppockCurve = indMomentum.NewCoppockCurveWithPeriods[float64](1, 1, 1)

	snapshots := []*asset.Snapshot{
		{Close: 10},
		{Close: 11},
		{Close: 11},
		{Close: 10},
	}

	actions := cc.ComputeWithContext(ctx, helper.SliceToChan(snapshots))
	actual := helper.ChanToSlice(actions)

	expected := []strategy.Action{
		strategy.Hold, // Idle period (index 0)
		strategy.Buy,  // Coppock > 0 (index 1)
		strategy.Hold, // Coppock == 0 (index 2)
		strategy.Sell, // Coppock < 0 (index 3)
	}

	if len(actual) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(actual))
	}

	for i, v := range actual {
		if v != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], v)
		}
	}
}
