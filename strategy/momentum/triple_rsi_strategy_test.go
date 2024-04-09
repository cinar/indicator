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

func TestTripleRsiStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/triple_rsi_strategy.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	tripleRsi := momentum.NewTripleRsiStrategy()
	actions, outcomes := strategy.ComputeWithOutcome(tripleRsi, snapshots)

	outcomes = helper.RoundDigits(outcomes, 2)

	err = strategy.CheckResults(results, actions, outcomes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTripleRsiStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	tripleRsi := momentum.NewTripleRsiStrategy()

	report := tripleRsi.Report(snapshots)

	fileName := "triple_rsi_strategy.html"
	defer os.Remove(fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
