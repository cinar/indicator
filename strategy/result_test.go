// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package strategy_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestCheckResults(t *testing.T) {
	results := helper.Duplicate(helper.SliceToChan([]*strategy.Result{
		{strategy.Buy, 0},
		{strategy.Hold, 1},
		{strategy.Hold, 2},
		{strategy.Sell, 3},
	}), 3)

	actions := helper.Map(results[1], func(r *strategy.Result) strategy.Action { return r.Action })
	outcomes := helper.Map(results[2], func(r *strategy.Result) float64 { return r.Outcome })

	err := strategy.CheckResults(results[0], actions, outcomes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckResultsActionsEnd(t *testing.T) {
	results := helper.Duplicate(helper.SliceToChan([]*strategy.Result{
		{strategy.Buy, 0},
		{strategy.Hold, 1},
		{strategy.Hold, 2},
		{strategy.Sell, 3},
	}), 3)

	actions := helper.Map(results[1], func(r *strategy.Result) strategy.Action { return r.Action })
	outcomes := helper.Map(results[2], func(r *strategy.Result) float64 { return r.Outcome })

	actions = helper.First(actions, 2)

	err := strategy.CheckResults(results[0], actions, outcomes)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckResultsOutcomesEnd(t *testing.T) {
	results := helper.Duplicate(helper.SliceToChan([]*strategy.Result{
		{strategy.Buy, 0},
		{strategy.Hold, 1},
		{strategy.Hold, 2},
		{strategy.Sell, 3},
	}), 3)

	actions := helper.Map(results[1], func(r *strategy.Result) strategy.Action { return r.Action })
	outcomes := helper.Map(results[2], func(r *strategy.Result) float64 { return r.Outcome })

	outcomes = helper.First(outcomes, 2)

	err := strategy.CheckResults(results[0], actions, outcomes)
	if err == nil {
		t.Fatal("expected error")
	}
}
