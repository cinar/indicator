// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestOrStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/or.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	or := strategy.NewOrStrategy("Or Strategy")
	or.Strategies = append(or.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	actual := or.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	or := strategy.NewOrStrategy("Or Strategy")
	or.Strategies = append(or.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	report := or.Report(snapshots)

	fileName := "or.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
