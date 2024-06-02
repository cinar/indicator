// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultApoFastPeriod is the default APO fast period of 14.
	DefaultApoFastPeriod = 14

	// DefaultApoFastSmoothing is the default APO fast smoothing.
	DefaultApoFastSmoothing = DefaultEmaSmoothing

	// DefaultApoSlowPeriod is the default APO slow period of 30.
	DefaultApoSlowPeriod = 30

	// DefaultApoSlowSmoothing is the default APO slow smoothing.
	DefaultApoSlowSmoothing = DefaultEmaSmoothing
)

// Apo represents the configuration parameters for calculating the
// Absolute Price Oscillator (APO). An APO value crossing above
// zero suggests a bullish trend, while crossing below zero
// indicates a bearish trend. Positive APO values signify
// an upward trend, while negative values signify a
// downward trend.
//
//	Fast = Ema(values, fastPeriod)
//	Slow = Ema(values, slowPeriod)
//	APO = Fast - Slow
//
// Example:
//
//	apo := trend.NewApo[float64]()
//	apo.FastPeriod = 12
//	apo.SlowPeriod = 26
//
//	result := apo.Compute(c)
type Apo[T helper.Number] struct {
	// Fast period.
	FastPeriod int

	// Fast smoothing.
	FastSmoothing T

	// Slow period.
	SlowPeriod int

	// Slow smoothing.
	SlowSmoothing T
}

// NewApo function initializes a new APO instance
// with the default parameters.
func NewApo[T helper.Number]() *Apo[T] {
	return &Apo[T]{
		FastPeriod:    DefaultApoFastPeriod,
		FastSmoothing: DefaultApoFastSmoothing,
		SlowPeriod:    DefaultApoSlowPeriod,
		SlowSmoothing: DefaultApoSlowSmoothing,
	}
}

// Compute function takes a channel of numbers and computes the APO
// over the specified period.
func (apo *Apo[T]) Compute(c <-chan T) <-chan T {
	c = helper.Buffered(c, apo.SlowPeriod)
	cs := helper.Duplicate(c, 2)

	fastEma := NewEma[T]()
	fastEma.Period = apo.FastPeriod
	cs[0] = fastEma.Compute(cs[0])

	slowEma := NewEma[T]()
	slowEma.Period = apo.SlowPeriod
	cs[1] = slowEma.Compute(cs[1])

	return helper.Subtract(cs[0], cs[1])
}
