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
//	Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
//	Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
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
	// The highest high and lowest low are looked for in the last Period + 1
	// bars, the current one included.
	sinceLastHigh := barsSinceLatest(ctx, high, a.Period+1, func(n, best T) bool { return n > best })
	sinceLastLow := barsSinceLatest(ctx, low, a.Period+1, func(n, best T) bool { return n < best })

	sinceLastHigh = helper.SkipWithContext(ctx, sinceLastHigh, a.Period)
	sinceLastLow = helper.SkipWithContext(ctx, sinceLastLow, a.Period)

	// Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
	aroonUp := helper.MultiplyByWithContext(ctx, sinceLastHigh, -1)
	aroonUp = helper.IncrementByWithContext(ctx, aroonUp, T(a.Period))
	aroonUp = helper.DivideByWithContext(ctx, aroonUp, T(a.Period))
	aroonUp = helper.MultiplyByWithContext(ctx, aroonUp, 100)
	aroonUp = helper.RoundDigitsWithContext(ctx, aroonUp, 0)

	// Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
	aroonDown := helper.MultiplyByWithContext(ctx, sinceLastLow, -1)
	aroonDown = helper.IncrementByWithContext(ctx, aroonDown, T(a.Period))
	aroonDown = helper.DivideByWithContext(ctx, aroonDown, T(a.Period))
	aroonDown = helper.MultiplyByWithContext(ctx, aroonDown, 100)
	aroonDown = helper.RoundDigitsWithContext(ctx, aroonDown, 0)

	return aroonUp, aroonDown
}

// IdlePeriod is the initial period that Aroon won't yield any results.
func (a *Aroon[T]) IdlePeriod() int {
	return a.Period
}

// String is the string representation of the Aroon.
func (a *Aroon[T]) String() string {
	return fmt.Sprintf("AROON(%d)", a.Period)
}

// barsSinceLatest returns how many bars ago the best value in each window of
// w bars was, where better reports whether n beats the best so far. On ties
// the most recent bar wins, as in TA-Lib.
func barsSinceLatest[T helper.Float](ctx context.Context, c <-chan T, w int, better func(n, best T) bool) <-chan T {
	return helper.WindowWithContext(ctx, c, func(s []T, i int) T {
		var best T
		since, k := 0, 0

		// Walk from the newest bar to the oldest.
		helper.SlicesReverse(s, i, func(n T) bool {
			if k == 0 || better(n, best) {
				best, since = n, k
			}
			k++
			return true
		})

		return T(since)
	}, w)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *Aroon[T]) Compute(high, low <-chan T) (<-chan T, <-chan T) {
	return a.ComputeWithContext(context.Background(), high, low)
}
