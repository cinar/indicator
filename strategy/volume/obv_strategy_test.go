// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/volume"
)

func TestObvStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/obv_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	obv := volume.NewObvStrategy()
	actual := obv.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestObvStrategyHold(t *testing.T) {
	// A mock setup where OBV equals SMA to test the Hold branch in Compute.
	// SMA of [10, 10] is 10. So if current OBV is 10, it's a Hold.
	snapshots := helper.SliceToChan([]*asset.Snapshot{
		{Close: 10, Volume: 10},
		{Close: 10, Volume: 0},
	})

	// 2-period SMA of OBV
	s := volume.NewObvStrategyWith(2)
	actions := s.Compute(snapshots)

	// Shift(actions, 1, Hold) -> first is Hold.
	// Second: OBV is 10, SMA is 10. OBV == SMA -> Hold.
	expected := []strategy.Action{strategy.Hold, strategy.Hold}

	err := helper.CheckEquals(actions, helper.SliceToChan(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestObvStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	obv := volume.NewObvStrategy()

	report := obv.Report(snapshots)

	fileName := "obv_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
