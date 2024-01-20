// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultEmaPeriod is the default EMA period of 20.
	DefaultEmaPeriod = 20

	// DefaultEmaSmoothing is the default EMA smooting of 2.
	DefaultEmaSmoothing = 2
)

// Ema represents the parameters for calculating the Exponential Moving Average.
//
// Example:
//
//	ema := trend.NewEma[float64]()
//	ema.Period = 10
//
//	result := ema.Compute(c)
type Ema[T helper.Number] struct {
	// Time period.
	Period int

	// Smoothing constant.
	Smoothing T
}

// NewEma function initializes a new EMA instance with the default parameters.
func NewEma[T helper.Number]() *Ema[T] {
	return &Ema[T]{
		Period:    DefaultEmaPeriod,
		Smoothing: DefaultEmaSmoothing,
	}
}

// NewEmaWithPeriod function initializes a new EMA instance with the given period.
func NewEmaWithPeriod[T helper.Number](period int) *Ema[T] {
	ema := NewEma[T]()
	ema.Period = period

	return ema
}

// Compute function takes a channel of numbers and computes the EMA over the specified period.
func (ema *Ema[T]) Compute(c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Initial EMA value is the SMA.
		sma := NewSma[T]()
		sma.Period = ema.Period

		before := <-sma.Compute(helper.Head(c, ema.Period))
		result <- before

		multiplier := ema.Smoothing / T(ema.Period+1)

		for n := range c {
			before = (n-before)*multiplier + before
			result <- before
		}
	}()

	return result
}

// IdlePeriod is the initial period that EMA yield any results.
func (ema *Ema[T]) IdlePeriod() int {
	return ema.Period - 1
}
