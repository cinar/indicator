// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestAndStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/and.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	and := strategy.NewAndStrategy("And Strategy")
	and.Strategies = append(and.Strategies, strategy.NewBuyAndHoldStrategy(), strategy.NewBuyAndHoldStrategy())

	actual := and.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAndStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	and := strategy.NewAndStrategy("And Strategy")
	and.Strategies = append(and.Strategies, strategy.NewBuyAndHoldStrategy(), strategy.NewBuyAndHoldStrategy())

	report := and.Report(snapshots)

	fileName := "and.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllAndStrategies(t *testing.T) {
	strategies := []strategy.Strategy{
		strategy.NewBuyAndHoldStrategy(),
		strategy.NewMajorityStrategyWith("", []strategy.Strategy{
			strategy.NewBuyAndHoldStrategy(),
		}),
	}

	allAndStrategies := strategy.AllAndStrategies(strategies)

	expected := len(strategies)*len(strategies) - len(strategies)
	actual := len(allAndStrategies)

	if actual != expected {
		t.Fatalf("actual=%d expected=%d", actual, expected)
	}
}
