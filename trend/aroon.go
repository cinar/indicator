// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultAroonPeriod is the default Aroon period of 25.
	DefaultAroonPeriod = 25
)

// Aroon represent the configuration for calculating the Aroon indicator. It is
// a technical analysis tool that gauges trend direction and strength in asset
// prices. It comprises two lines: Aroon Up and Aroon Down. Aroon Up measures
// uptrend strength, while Aroon Down measures downtrend strength. When Aroon
// Up exceeds Aroon Down, it suggests a bullish trend; when Aroon Down
// surpasses Aroon Up, it indicates a bearish trend.
//
//	Aroon Up = ((Period - Bars Since Highest High) / Period) * 100
//	Aroon Down = ((Period - Bars Since Lowest Low) / Period) * 100
//
// The highest high and lowest low are looked for in the last Period + 1
// bars, the current one included. When the extreme occurs more than once,
// the most recent one counts, as in TA-Lib. A Period below 1 produces no
// output.
//
// Example:
//
//	aroon := trend.NewAroon[float64]()
//	aroon.Period = 25
//
//	result := aroon.Compute(c)
type Aroon[T helper.Float] struct {
	// Period is the period to use.
	Period int
}

// NewAroon function initializes a new Aroon instance
// with the default parameters.
func NewAroon[T helper.Float]() *Aroon[T] {
	return &Aroon[T]{
		Period: DefaultAroonPeriod,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Aroon
// over the specified period.
func (a *Aroon[T]) ComputeWithContext(ctx context.Context, high, low <-chan T) (<-chan T, <-chan T) {
	if a.Period < 1 {
		// A window of size 0 produces no output.
		return helper.MaxSinceWithContext(ctx, high, 0), helper.MinSinceWithContext(ctx, low, 0)
	}

	// The highest high and lowest low are looked for in the last Period + 1
	// bars, the current one included.
	sinceLastHigh := helper.MaxSinceWithContext(ctx, high, a.Period+1)
	sinceLastLow := helper.MinSinceWithContext(ctx, low, a.Period+1)

	sinceLastHigh = helper.SkipWithContext(ctx, sinceLastHigh, a.Period)
	sinceLastLow = helper.SkipWithContext(ctx, sinceLastLow, a.Period)

	return a.fromBarsSince(ctx, sinceLastHigh), a.fromBarsSince(ctx, sinceLastLow)
}

// fromBarsSince converts the bars since an extreme to an Aroon line. It is
// computed as (Period - Bars Since) * (100 / Period), in the same order as
// TA-Lib, so that the floating-point results are the same.
func (a *Aroon[T]) fromBarsSince(ctx context.Context, since <-chan T) <-chan T {
	line := helper.MultiplyByWithContext(ctx, since, -1)
	line = helper.IncrementByWithContext(ctx, line, T(a.Period))
	return helper.MultiplyByWithContext(ctx, line, 100/T(a.Period))
}

// IdlePeriod is the initial period that Aroon won't yield any results.
func (a *Aroon[T]) IdlePeriod() int {
	return a.Period
}

// String is the string representation of the Aroon.
func (a *Aroon[T]) String() string {
	return fmt.Sprintf("AROON(%d)", a.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *Aroon[T]) Compute(high, low <-chan T) (<-chan T, <-chan T) {
	return a.ComputeWithContext(context.Background(), high, low)
}
