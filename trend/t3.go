// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"
	"math"

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
//	c2 = 3a^2
//	c3 = -3a
//	c4 = a^3
//	a = volume factor
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

// Compute function takes a channel of numbers and computes the T3 Moving Average.
func (t *T3[T]) Compute(closings <-chan T) <-chan T {
	// Chain 6 EMAs
	ema1 := t.ema1.Compute(closings)
	ema2 := t.ema2.Compute(ema1)
	ema3 := t.ema3.Compute(ema2)
	ema4 := t.ema4.Compute(ema3)
	ema5 := t.ema5.Compute(ema4)
	ema6 := t.ema6.Compute(ema5)

	// Calculate coefficients based on volume factor
	a := float64(t.VolumeFactor)
	c1 := -math.Pow(a, 3)
	c2 := 3 * math.Pow(a, 2)
	c3 := -3 * a
	c4 := math.Pow(a, 3)

	// T3 = c1*EMA6 + c2*EMA6(EMA6) + c3*EMA6(EMA6(EMA6)) + c4*EMA6(EMA6(EMA6(EMA6)))
	// Which is: c1*ema6 + c2*ema5 + c3*ema4 + c4*ema3
	result := helper.Add(
		helper.Add(
			helper.MultiplyBy(ema6, T(c1)),
			helper.MultiplyBy(ema5, T(c2)),
		),
		helper.Add(
			helper.MultiplyBy(ema4, T(c3)),
			helper.MultiplyBy(ema3, T(c4)),
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
