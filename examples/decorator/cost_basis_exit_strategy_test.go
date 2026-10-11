// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/examples/decorator"
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
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

// fixedActionsStrategy is a test double that ignores the snapshot values and emits a
// predetermined sequence of actions, one per snapshot received. Any snapshot beyond the
// end of the configured sequence results in Hold.
type fixedActionsStrategy struct {
	actions []strategy.Action
}

func (s *fixedActionsStrategy) Name() string {
	return "fixedActionsStrategy"
}

func (s *fixedActionsStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	actions := make(chan strategy.Action)

	go func() {
		defer close(actions)

		i := 0
		for range snapshots {
			if i < len(s.actions) {
				actions <- s.actions[i]
			} else {
				actions <- strategy.Hold
			}
			i++
		}
	}()

	return actions
}

func (s *fixedActionsStrategy) Report(_ <-chan *asset.Snapshot) *helper.Report {
	return &helper.Report{}
}

// TestCostBasisExitStrategyBreakEven confirms that a Sell signal whose closing price is
// exactly equal to the earlier Buy's closing price (break even) is honored as a Sell,
// matching the type's documented "at or above the original purchase price" behavior.
// Prior to the fix, the strict less-than comparison caused this case to emit Hold instead.
func TestCostBasisExitStrategyBreakEven(t *testing.T) {
	innerStrategy := &fixedActionsStrategy{
		actions: []strategy.Action{strategy.Buy, strategy.Hold, strategy.Sell},
	}

	snapshots := helper.SliceToChan([]*asset.Snapshot{
		{Close: 10},
		{Close: 8},
		{Close: 10},
	})

	costBasisExitStrategy := decorator.NewCostBasisExitStrategy(innerStrategy)

	actual := costBasisExitStrategy.Compute(snapshots)
	expected := helper.SliceToChan([]strategy.Action{strategy.Buy, strategy.Hold, strategy.Sell})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
