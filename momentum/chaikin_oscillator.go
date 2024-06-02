// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volume"
)

const (
	// DefaultChaikinOscillatorShortPeriod is the default short period for the Chaikin Oscillator.
	DefaultChaikinOscillatorShortPeriod = 3

	// DefaultChaikinOscillatorLongPeriod is the default long period for the Chaikin Oscillator.
	DefaultChaikinOscillatorLongPeriod = 10
)

// ChaikinOscillator represents the configuration parameter for calculating the Chaikin Oscillator. It measures
// the momentum of the Accumulation/Distribution (A/D) using the Moving Average Convergence Divergence (MACD)
// formula. It takes the difference between fast and slow periods EMA of the A/D. Cross above the A/D line
// indicates bullish.
//
//	CO = Ema(fastPeriod, AD) - Ema(slowPeriod, AD)
//
// Example:
//
//	co := momentum.ChaikinOscillator[float64]()
//	values := co.Compute(lows, highs)
type ChaikinOscillator[T helper.Number] struct {
	// Ad is the Accumulation/Distribution (A/D) instance.
	Ad *volume.Ad[T]

	// ShortEma is the SMA for the short period.
	ShortEma *trend.Ema[T]

	// LongEma is the SMA for the long period.
	LongEma *trend.Ema[T]
}

// NewChaikinOscillator function initializes a new Chaikin Oscillator instance.
func NewChaikinOscillator[T helper.Number]() *ChaikinOscillator[T] {
	return &ChaikinOscillator[T]{
		Ad:       volume.NewAd[T](),
		ShortEma: trend.NewEmaWithPeriod[T](DefaultChaikinOscillatorShortPeriod),
		LongEma:  trend.NewEmaWithPeriod[T](DefaultChaikinOscillatorLongPeriod),
	}
}

// Compute function takes a channel of numbers and computes the Chaikin Oscillator.
func (c *ChaikinOscillator[T]) Compute(highs, lows, closings, volumes <-chan T) (<-chan T, <-chan T) {
	adSplice := helper.Duplicate(
		c.Ad.Compute(highs, lows, closings, volumes),
		3,
	)

	shortEma := c.ShortEma.Compute(adSplice[0])
	longEma := c.LongEma.Compute(adSplice[1])

	shortEma = helper.Skip(shortEma, c.LongEma.IdlePeriod()-c.ShortEma.IdlePeriod())

	co := helper.Subtract(shortEma, longEma)
	adSplice[2] = helper.Skip(adSplice[2], c.LongEma.IdlePeriod())

	return co, adSplice[2]
}

// IdlePeriod is the initial period that Chaikin Oscillator won't yield any results.
func (c *ChaikinOscillator[T]) IdlePeriod() int {
	return c.LongEma.IdlePeriod()
}
