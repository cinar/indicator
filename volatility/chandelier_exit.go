// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultChandelierExitPeriod is the default period for the Chandelier Exit.
	DefaultChandelierExitPeriod = 22

	// DefaultChandelierExitMultiplier is the default multiplier for the Chandelier Exit.
	DefaultChandelierExitMultiplier = 3
)

// ChandelierExit represents the configuration parameters for calculating the Chandelier Exit.
// It sets a trailing stop-loss based on the Average True Value (ATR).
//
//	Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
//	Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
//
// Example:
//
//	ce := volatility.NewChandelierExit[float64]()
//	ceLong, ceShort := ce.Compute(highs, lows, closings)
type ChandelierExit[T helper.Number] struct {
	// Period is time period.
	Period int

	// Multiplier is for sensitivity.
	Multiplier T
}

// NewChandelierExit function initializes a new Chandelier Exit instance with the default parameters.
func NewChandelierExit[T helper.Number]() *ChandelierExit[T] {
	return &ChandelierExit[T]{
		Period:     DefaultChandelierExitPeriod,
		Multiplier: DefaultChandelierExitMultiplier,
	}
}

// Compute function takes a channel of numbers and computes the Chandelier Exit over the specified period.
func (c *ChandelierExit[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T) {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)

	movingMax := trend.NewMovingMaxWithPeriod[T](c.Period)
	movingMin := trend.NewMovingMinWithPeriod[T](c.Period)

	atr := NewAtrWithPeriod[T](c.Period)

	maxHighs := helper.Skip(
		movingMax.Compute(highsSplice[0]),
		atr.IdlePeriod()-movingMax.IdlePeriod(),
	)

	minLows := helper.Skip(
		movingMin.Compute(lowsSplice[0]),
		atr.IdlePeriod()-movingMin.IdlePeriod(),
	)

	atr3Splice := helper.Duplicate(
		helper.MultiplyBy(
			atr.Compute(highsSplice[1], lowsSplice[1], closings),
			c.Multiplier,
		),
		2,
	)

	ceLong := helper.Subtract(maxHighs, atr3Splice[0])
	ceShort := helper.Add(minLows, atr3Splice[1])

	return ceLong, ceShort
}

// IdlePeriod is the initial period that Chandelier Exit won't yield any results.
func (c *ChandelierExit[T]) IdlePeriod() int {
	return c.Period
}
