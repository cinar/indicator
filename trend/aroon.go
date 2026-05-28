// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

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
type Aroon[T helper.Number] struct {
	// Period is the period to use.
	Period int
}

// NewAroon function initializes a new Aroon instance
// with the default parameters.
func NewAroon[T helper.Number]() *Aroon[T] {
	return &Aroon[T]{
		Period: DefaultAroonPeriod,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Aroon
// over the specified period.
func (a *Aroon[T]) ComputeWithContext(ctx context.Context, high, low <-chan T) (<-chan T, <-chan T) {
	movingMax := NewMovingMaxWithPeriod[T](a.Period)
	movingMin := NewMovingMinWithPeriod[T](a.Period)

	sinceLastHigh := helper.MaxSince(movingMax.ComputeWithContext(ctx, high), a.Period)
	sinceLastLow := helper.MinSince(movingMin.ComputeWithContext(ctx, low), a.Period)

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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *Aroon[T]) Compute(high, low <-chan T) (<-chan T, <-chan T) {
	return a.ComputeWithContext(context.Background(), high, low)
}
