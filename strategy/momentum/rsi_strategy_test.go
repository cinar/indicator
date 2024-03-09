// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package momentum_test

import (
	"os"
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/momentum"
)

func TestRsiStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/rsi_strategy.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	rsi := momentum.NewRsiStrategy()
	actions, outcomes := strategy.ComputeWithOutcome(rsi, snapshots)

	outcomes = helper.RoundDigits(outcomes, 2)

	err = strategy.CheckResults(results, actions, outcomes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRsiStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	rsi := momentum.NewRsiStrategy()

	report := rsi.Report(snapshots)

	fileName := "rsi_strategy.html"
	defer os.Remove(fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
