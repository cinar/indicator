// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultT3Period is the default period for the T3 Moving Average.
	DefaultT3Period = 5

	// DefaultT3VolumeFactor is the default volume factor for T3.
	DefaultT3VolumeFactor = 0.7
)

// T3 represents the configuration parameters for calculating the
// Tillson T3 Moving Average. The T3 is a smooth moving average
// that chains multiple EMAs together with a volume factor for
// improved responsiveness.
//
//	T3 = c1*EMA6 + c2*EMA6(EMA6) + c3*EMA6(EMA6(EMA6)) + c4*EMA6(EMA6(EMA6(EMA6)))
//
// where:
//
//	c1 = -a^3
//	c2 = 3a^2 + 3a^3
//	c3 = -6a^2 - 3a - 3a^3
//	c4 = 1 + 3a + 3a^2 + a^3
//	a = volume factor
//
// The coefficients sum to 1 for any a, the way a weighted moving average's
// must, so T3 reproduces a constant input exactly.
//
// Example:
//
//	t3 := trend.NewT3[float64]()
//	result := t3.Compute(closings)
type T3[T helper.Float] struct {
	// Period is the period for the EMA calculations.
	Period int

	// VolumeFactor is the volume factor for the T3 calculation.
	VolumeFactor T

	// ema1 through ema6 are the EMA instances for chaining.
	ema1, ema2, ema3, ema4, ema5, ema6 *Ema[T]
}

// NewT3 function initializes a new T3 instance.
func NewT3[T helper.Float]() *T3[T] {
	return NewT3WithPeriodAndFactor[T](DefaultT3Period, DefaultT3VolumeFactor)
}

// NewT3WithPeriodAndFactor function initializes a new T3 instance with
// specified period and volume factor.
func NewT3WithPeriodAndFactor[T helper.Float](period int, volumeFactor float64) *T3[T] {
	t := &T3[T]{
		Period:       period,
		VolumeFactor: T(volumeFactor),
	}

	// Create 6 chained EMA instances
	t.ema1 = NewEmaWithPeriod[T](period)
	t.ema2 = NewEmaWithPeriod[T](period)
	t.ema3 = NewEmaWithPeriod[T](period)
	t.ema4 = NewEmaWithPeriod[T](period)
	t.ema5 = NewEmaWithPeriod[T](period)
	t.ema6 = NewEmaWithPeriod[T](period)

	return t
}

// ComputeWithContext function takes a channel of numbers and computes the T3 Moving Average.
func (t *T3[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	// Chain 6 EMAs. ema3, ema4, and ema5 each feed the next EMA in the
	// chain *and* are used directly below in the weighted sum, so each
	// needs its own duplicated copy for the second use -- a single
	// channel only has one consumer's worth of values to give out.
	ema1 := t.ema1.ComputeWithContext(ctx, closings)
	ema2 := t.ema2.ComputeWithContext(ctx, ema1)

	ema3Splice := helper.DuplicateWithContext(ctx, t.ema3.ComputeWithContext(ctx, ema2), 2)
	ema4Splice := helper.DuplicateWithContext(ctx, t.ema4.ComputeWithContext(ctx, ema3Splice[0]), 2)
	ema5Splice := helper.DuplicateWithContext(ctx, t.ema5.ComputeWithContext(ctx, ema4Splice[0]), 2)
	ema6 := t.ema6.ComputeWithContext(ctx, ema5Splice[0])

	// Each EMA only starts yielding once it has Period-1 more inputs than
	// it was given, so ema3/ema4/ema5's channels run ahead of ema6's by
	// 3/2/1 EMA delays respectively. Skip that many off each so that,
	// read in lockstep with ema6, they land on the same closing.
	idle := t.Period - 1
	ema3Aligned := helper.SkipWithContext(ctx, ema3Splice[1], 3*idle)
	ema4Aligned := helper.SkipWithContext(ctx, ema4Splice[1], 2*idle)
	ema5Aligned := helper.SkipWithContext(ctx, ema5Splice[1], idle)

	// Calculate coefficients based on volume factor. These sum to 1 for
	// any a (verify: -a^3 + 3a^2+3a^3 -6a^2-3a-3a^3 + 1+3a+3a^2+a^3 = 1),
	// as required of a weighted moving average.
	a := float64(t.VolumeFactor)
	a2 := a * a
	a3 := a2 * a
	c1 := -a3
	c2 := 3*a2 + 3*a3
	c3 := -6*a2 - 3*a - 3*a3
	c4 := 1 + 3*a + 3*a2 + a3

	// T3 = c1*EMA6 + c2*EMA6(EMA6) + c3*EMA6(EMA6(EMA6)) + c4*EMA6(EMA6(EMA6(EMA6)))
	// Which is: c1*ema6 + c2*ema5 + c3*ema4 + c4*ema3
	result := helper.AddWithContext(ctx, helper.AddWithContext(ctx, helper.MultiplyByWithContext(ctx, ema6, T(c1)),
		helper.MultiplyByWithContext(ctx, ema5Aligned, T(c2)),
	),
		helper.AddWithContext(ctx, helper.MultiplyByWithContext(ctx, ema4Aligned, T(c3)),
			helper.MultiplyByWithContext(ctx, ema3Aligned, T(c4)),
		),
	)

	return result
}

// IdlePeriod is the initial period that T3 won't yield any results.
func (t *T3[T]) IdlePeriod() int {
	// Each EMA has a delay of Period-1, and we chain 6 EMAs
	// Total delay = 6 * (Period - 1)
	return 6 * (t.Period - 1)
}

// String is the string representation of the T3.
func (t *T3[T]) String() string {
	return fmt.Sprintf("T3(%d, %.1f)", t.Period, float64(t.VolumeFactor))
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (t *T3[T]) Compute(closings <-chan T) <-chan T {
	return t.ComputeWithContext(context.Background(), closings)
}
