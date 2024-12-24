// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
const (
	// DefaultAccelerationBandsPeriod is the default period for the Acceleration Bands.
	DefaultAccelerationBandsPeriod = 20
)

// AccelerationBands represents the configuration parameters for calculating the Acceleration Bands.
//
//	Upper Band = SMA(High * (1 + 4 * (High - Low) / (High + Low)))
//	Middle Band = SMA(Closing)
//	Lower Band = SMA(Low * (1 - 4 * (High - Low) / (High + Low)))
//
// Example:
//
//	accelerationBands := NewAccelerationBands[float64]()
//	accelerationBands.Compute(values)
type AccelerationBands[T helper.Number] struct {
	// Time period.
	Period int
}

// NewAccelerationBands function initializes a new Acceleration Bands instance with the default parameters.
func NewAccelerationBands[T helper.Number]() *AccelerationBands[T] {
	return &AccelerationBands[T]{
		Period: DefaultAccelerationBandsPeriod,
	}
}

// Compute function takes a channel of numbers and computes the Acceleration Bands over the specified period.
func (a *AccelerationBands[T]) Compute(high, low, closing <-chan T) (<-chan T, <-chan T, <-chan T) {
	highs := helper.Duplicate(high, 3)
	lows := helper.Duplicate(low, 3)

	ks := helper.Duplicate(
		helper.Divide(
			helper.Subtract(highs[0], lows[0]),
			helper.Add(highs[1], lows[1]),
		),
		2,
	)

	sma := trend.NewSmaWithPeriod[T](a.Period)

	upper := sma.Compute(
		helper.Multiply(
			highs[2],
			helper.IncrementBy(
				helper.MultiplyBy(
					ks[0],
					4,
				),
				1,
			),
		),
	)

	middle := sma.Compute(closing)

	lower := sma.Compute(
		helper.Multiply(
			lows[2],
			helper.IncrementBy(
				helper.MultiplyBy(
					ks[1],
					-4,
				),
				1,
			),
		),
	)

	return upper, middle, lower
}

// IdlePeriod is the initial period that Acceleration Bands won't yield any results.
func (a *AccelerationBands[T]) IdlePeriod() int {
	return a.Period - 1
}
