// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/decorator"
	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestInverseStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/inverse_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	innerStrategy := trend.NewMacdStrategy()
	inverseStrategy := decorator.NewInverseStrategy(innerStrategy)

	actual := inverseStrategy.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInverseStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	innerStrategy := trend.NewMacdStrategy()
	inverseStrategy := decorator.NewInverseStrategy(innerStrategy)

	report := inverseStrategy.Report(snapshots)

	fileName := "inverse_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
