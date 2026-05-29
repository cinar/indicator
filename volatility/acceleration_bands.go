// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"context"

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

// ComputeWithContext function takes a channel of numbers and computes the Acceleration Bands over the specified period.
func (a *AccelerationBands[T]) ComputeWithContext(ctx context.Context, high, low, closing <-chan T) (<-chan T, <-chan T, <-chan T) {
	highs := helper.DuplicateWithContext(ctx, high, 3)
	lows := helper.DuplicateWithContext(ctx, low, 3)

	ks := helper.DuplicateWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, highs[0], lows[0]),
		helper.AddWithContext(ctx, highs[1], lows[1]),
	),
		2,
	)

	sma := trend.NewSmaWithPeriod[T](a.Period)

	upper := sma.ComputeWithContext(ctx, helper.MultiplyWithContext(ctx, highs[2],
		helper.IncrementByWithContext(ctx, helper.MultiplyByWithContext(ctx, ks[0],
			4,
		),
			1,
		),
	),
	)

	middle := sma.ComputeWithContext(ctx, closing)

	lower := sma.ComputeWithContext(ctx, helper.MultiplyWithContext(ctx, lows[2],
		helper.IncrementByWithContext(ctx, helper.MultiplyByWithContext(ctx, ks[1],
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *AccelerationBands[T]) Compute(high, low, closing <-chan T) (<-chan T, <-chan T, <-chan T) {
	return a.ComputeWithContext(context.Background(), high, low, closing)
}
