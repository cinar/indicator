// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

// ExecutionTiming represents when a simulated trade executes relative to the bar its action was computed from.
type ExecutionTiming int

const (
	// AtClose executes at the same bar's closing price. This is the existing default behavior used by
	// OutcomeWithContext and ComputeWithOutcomeWithContext, and assumes the trade fills at the very
	// close that produced the signal.
	AtClose ExecutionTiming = iota

	// NextOpen executes at the opening price of the bar immediately following the signal. This is a
	// more realistic assumption for strategies that can only act after a bar has fully closed.
	NextOpen

	// NextClose executes at the closing price of the bar immediately following the signal.
	NextClose
)

// String returns the string representation of the ExecutionTiming.
func (e ExecutionTiming) String() string {
	switch e {
	case NextOpen:
		return "NextOpen"

	case NextClose:
		return "NextClose"

	default:
		return "AtClose"
	}
}
