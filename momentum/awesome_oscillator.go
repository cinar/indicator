// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultAwesomeOscillatorShortPeriod is the default short period for the Awesome Oscillator (AO).
	DefaultAwesomeOscillatorShortPeriod = 5

	// DefaultAwesomeOscillatorLongPeriod is the default long period for the Awesome Oscillator (AO).
	DefaultAwesomeOscillatorLongPeriod = 34
)

// AwesomeOscillator represents the configuration parameter for calculating the Awesome Oscillator (AO). It gauges
// market momentum by comparing short-term price action (5-period average) against long-term trends (34-period
// average). Its value around a zero line reflects bullishness above and bearishness below. Crossings of the
// zero line can signal potential trend reversals. Traders use the AO to confirm existing trends, identify
// entry/exit points, and understand momentum shifts.
//
//	Median Price = ((Low + High) / 2).
//	AO = 5-Period SMA - 34-Period SMA.
//
// Example:
//
//	ao := momentum.AwesomeOscillator[float64]()
//	values := ao.Compute(lows, highs)
type AwesomeOscillator[T helper.Number] struct {
	// ShortSma is the SMA for the short period.
	ShortSma *trend.Sma[T]

	// LongSma is the SMA for the long period.
	LongSma *trend.Sma[T]
}

// NewAwesomeOscillator function initializes a new Awesome Oscillator instance.
func NewAwesomeOscillator[T helper.Number]() *AwesomeOscillator[T] {
	return &AwesomeOscillator[T]{
		ShortSma: trend.NewSmaWithPeriod[T](DefaultAwesomeOscillatorShortPeriod),
		LongSma:  trend.NewSmaWithPeriod[T](DefaultAwesomeOscillatorLongPeriod),
	}
}

// Compute function takes a channel of numbers and computes the AwesomeOscillator.
func (a *AwesomeOscillator[T]) Compute(highs, lows <-chan T) <-chan T {
	medianSplice := helper.Duplicate(
		helper.DivideBy(
			helper.Add(highs, lows),
			2,
		),
		2,
	)

	shortSma := a.ShortSma.Compute(medianSplice[0])
	longSma := a.LongSma.Compute(medianSplice[1])

	shortSma = helper.Skip(shortSma, a.LongSma.IdlePeriod()-a.ShortSma.IdlePeriod())

	return helper.Subtract(
		shortSma,
		longSma,
	)
}

// IdlePeriod is the initial period that Awesome Oscillator won't yield any results.
func (a *AwesomeOscillator[T]) IdlePeriod() int {
	return a.LongSma.IdlePeriod()
}
