// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestAlligatorStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/alligator_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	alligator := trend.NewAlligatorStrategy()
	actual := alligator.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlligatorStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	alligator := trend.NewAlligatorStrategy()
	report := alligator.Report(snapshots)

	fileName := "alligator_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
