// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultUlcerIndexPeriod is the default period for the Ulcer Index.
	DefaultUlcerIndexPeriod = 14
)

// UlcerIndex represents the configuration parameters for calculating the Ulcer Index (UI).
// It measures downside risk. The index increases in value as the price moves farther away
// from a recent high and falls as the price rises to new highs.
//
//	High Closings = Max(period, Closings)
//	Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
//	Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
//	Ulcer Index = Sqrt(Squared Average)
//
// Example:
//
//	ui := volatility.NewUlcerIndex[float64]()
//	ui.Compute(closings)
type UlcerIndex[T helper.Number] struct {
	// Time period.
	Period int
}

// NewUlcerIndex function initializes a new Ulcer Index instance with the default parameters.
func NewUlcerIndex[T helper.Number]() *UlcerIndex[T] {
	return &UlcerIndex[T]{
		Period: DefaultUlcerIndexPeriod,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Ulcer Index over the specified period.
func (u *UlcerIndex[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)

	//	High Closings = Max(period, Closings)
	movingMax := trend.NewMovingMaxWithPeriod[T](u.Period)
	highsSplice := helper.DuplicateWithContext(ctx, movingMax.ComputeWithContext(ctx, closingsSplice[0]),
		2,
	)

	//	Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
	closingsSplice[1] = helper.SkipWithContext(ctx, closingsSplice[1], movingMax.Period-1)

	percentageDrawdown := helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, closingsSplice[1], highsSplice[0]),
		highsSplice[1],
	),
		100,
	)

	//	Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
	sma := trend.NewSmaWithPeriod[T](u.Period)
	squaredAverage := helper.PowWithContext(ctx, sma.ComputeWithContext(ctx, percentageDrawdown),
		2,
	)

	// Ulcer Index = Sqrt(Squared Average)
	ulcerIndex := helper.SqrtWithContext(ctx, squaredAverage)

	return ulcerIndex
}

// IdlePeriod is the initial period that Ulcer Index won't yield any results.
func (u *UlcerIndex[T]) IdlePeriod() int {
	return (u.Period - 1) * 2
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (u *UlcerIndex[T]) Compute(closings <-chan T) <-chan T {
	return u.ComputeWithContext(context.Background(), closings)
}
