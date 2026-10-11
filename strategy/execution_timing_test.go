// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

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

// executionTimingSnapshots returns a small, hand-picked 4-bar OHLC series where the
// opening and closing prices differ enough, bar over bar, that AtClose, NextOpen, and
// NextClose execution timings produce clearly distinct, independently verifiable outcomes.
func executionTimingSnapshots() <-chan *asset.Snapshot {
	return helper.SliceToChan([]*asset.Snapshot{
		{Open: 10, Close: 12},
		{Open: 13, Close: 11},
		{Open: 10, Close: 14},
		{Open: 15, Close: 13},
	})
}

// executionTimingActions is Buy on the first bar, Sell on the third bar, Hold otherwise.
func executionTimingActions() []strategy.Action {
	return []strategy.Action{strategy.Buy, strategy.Hold, strategy.Sell, strategy.Hold}
}

// drainActions consumes and discards an actions channel in the background. The actions channel
// returned alongside outcomes is one branch of an internal Duplicate; it must be drained
// concurrently with outcomes or the shared upstream pipeline will block.
func drainActions(actions <-chan strategy.Action) {
	go func() {
		for range actions {
		}
	}()
}

func TestComputeWithOutcomeAndTimingAtClose(t *testing.T) {
	s := &fixedActionsStrategy{actions: executionTimingActions()}

	actions, outcomes := strategy.ComputeWithOutcomeAndTiming(s, executionTimingSnapshots(), strategy.AtClose)
	drainActions(actions)

	// Hand computed: buy at close 12 (shares = 1/12), hold at close 11, sell at close 14
	// (balance = (1/12)*14), hold at close 13 (balance unchanged since fully in cash).
	expected := helper.SliceToChan([]float64{0, -0.08, 0.17, 0.17})

	actual := helper.RoundDigits(outcomes, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestComputeWithOutcomeAndTimingNextOpen(t *testing.T) {
	s := &fixedActionsStrategy{actions: executionTimingActions()}

	actions, outcomes := strategy.ComputeWithOutcomeAndTiming(s, executionTimingSnapshots(), strategy.NextOpen)
	drainActions(actions)

	// Hand computed: the Buy signal on bar 1 executes at bar 2's open (13), the Sell signal on
	// bar 3 executes at bar 4's open (15). The action on the last bar has no following bar and
	// is dropped, so this channel is one element shorter than the actions channel.
	//   buy at 13 (shares = 1/13) -> outcome 0
	//   hold, valued at open 10   -> outcome (1/13)*10 - 1 = -0.230769...
	//   sell at 15 (balance = (1/13)*15) -> outcome 0.153846...
	expected := helper.SliceToChan([]float64{0, -0.23, 0.15})

	actual := helper.RoundDigits(outcomes, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestComputeWithOutcomeAndTimingNextClose(t *testing.T) {
	s := &fixedActionsStrategy{actions: executionTimingActions()}

	actions, outcomes := strategy.ComputeWithOutcomeAndTiming(s, executionTimingSnapshots(), strategy.NextClose)
	drainActions(actions)

	// Hand computed: the Buy signal on bar 1 executes at bar 2's close (11), the Sell signal on
	// bar 3 executes at bar 4's close (13). As with NextOpen, the last bar's action is dropped
	// since there is no following bar.
	//   buy at 11 (shares = 1/11) -> outcome 0
	//   hold, valued at close 14  -> outcome (1/11)*14 - 1 = 0.272727...
	//   sell at 13 (balance = (1/11)*13) -> outcome 0.181818...
	expected := helper.SliceToChan([]float64{0, 0.27, 0.18})

	actual := helper.RoundDigits(outcomes, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// TestComputeWithOutcomeAndTimingAtCloseMatchesDefault confirms that ExecutionTiming AtClose is
// byte-for-byte identical to the existing default behavior of ComputeWithOutcome and
// ComputeWithOutcomeWithContext, i.e. adding ExecutionTiming introduces zero behavior change to
// the pre-existing default execution path.
func TestComputeWithOutcomeAndTimingAtCloseMatchesDefault(t *testing.T) {
	s := &fixedActionsStrategy{actions: executionTimingActions()}

	defaultActions, defaultOutcomes := strategy.ComputeWithOutcome(s, executionTimingSnapshots())
	drainActions(defaultActions)

	timingActions, timingOutcomes := strategy.ComputeWithOutcomeAndTiming(s, executionTimingSnapshots(), strategy.AtClose)
	drainActions(timingActions)

	err := helper.CheckEquals(timingOutcomes, defaultOutcomes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecutionTimingString(t *testing.T) {
	tests := []struct {
		timing   strategy.ExecutionTiming
		expected string
	}{
		{strategy.AtClose, "AtClose"},
		{strategy.NextOpen, "NextOpen"},
		{strategy.NextClose, "NextClose"},
	}

	for _, test := range tests {
		actual := test.timing.String()
		if actual != test.expected {
			t.Fatalf("expected %s, got %s", test.expected, actual)
		}
	}
}
