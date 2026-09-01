// Copyright (c) 2021-2026 The Indicator Authors.
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

func TestCostBasisExitStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/cost_basis_exit_strategy.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	innerStrategy := trend.NewAroonStrategy()
	costBasisExitStrategy := decorator.NewCostBasisExitStrategy(innerStrategy)

	actual := costBasisExitStrategy.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCostBasisExitStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	innerStrategy := trend.NewAroonStrategy()
	costBasisExitStrategy := decorator.NewCostBasisExitStrategy(innerStrategy)

	report := costBasisExitStrategy.Report(snapshots)

	fileName := "cost_basis_exit_strategy.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoLossStrategyAlias(t *testing.T) {
	innerStrategy := trend.NewAroonStrategy()
	s := decorator.NewNoLossStrategy(innerStrategy)
	if s == nil {
		t.Fatal("expected non-nil strategy")
	}
}
