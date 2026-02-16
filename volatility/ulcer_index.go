// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
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

// Compute function takes a channel of numbers and computes the Ulcer Index over the specified period.
func (u *UlcerIndex[T]) Compute(closings <-chan T) <-chan T {
	closingsSplice := helper.Duplicate(closings, 2)

	//	High Closings = Max(period, Closings)
	movingMax := trend.NewMovingMaxWithPeriod[T](u.Period)
	highsSplice := helper.Duplicate(
		movingMax.Compute(closingsSplice[0]),
		2,
	)

	//	Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
	closingsSplice[1] = helper.Skip(closingsSplice[1], movingMax.Period-1)

	percentageDrawdown := helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(closingsSplice[1], highsSplice[0]),
			highsSplice[1],
		),
		100,
	)

	//	Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
	sma := trend.NewSmaWithPeriod[T](u.Period)
	squaredAverage := helper.Pow(
		sma.Compute(percentageDrawdown),
		2,
	)

	// Ulcer Index = Sqrt(Squared Average)
	ulcerIndex := helper.Sqrt(squaredAverage)

	return ulcerIndex
}

// IdlePeriod is the initial period that Ulcer Index won't yield any results.
func (u *UlcerIndex[T]) IdlePeriod() int {
	return (u.Period - 1) * 2
}
