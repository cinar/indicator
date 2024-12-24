// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultBollingerBandsPeriod is the default period for the Bollinger Bands.
	DefaultBollingerBandsPeriod = 20
)

// BollingerBands represents the configuration parameters for calculating the Bollinger Bands. It is a technical
// analysis tool used to gauge a market's volatility and identify overbought and oversold conditions. Returns
// the upper band, the middle band, and the lower band.
//
//	Middle Band = 20-Period SMA.
//	Upper Band = 20-Period SMA + 2 (20-Period Std)
//	Lower Band = 20-Period SMA - 2 (20-Period Std)
//
// Example:
//
//	bollingerBands := NewBollingerBands[float64]()
//	bollingerBands.Compute(values)
type BollingerBands[T helper.Number] struct {
	// Time period.
	Period int
}

// NewBollingerBands function initializes a new Bollinger Bands instance with the default parameters.
func NewBollingerBands[T helper.Number]() *BollingerBands[T] {
	return NewBollingerBandsWithPeriod[T](DefaultBollingerBandsPeriod)
}

// NewBollingerBandsWithPeriod function initializes a new Bollinger Bands instance with the given period.
func NewBollingerBandsWithPeriod[T helper.Number](period int) *BollingerBands[T] {
	return &BollingerBands[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the Bollinger Bands over the specified period.
func (b *BollingerBands[T]) Compute(c <-chan T) (<-chan T, <-chan T, <-chan T) {
	cs := helper.Duplicate(c, 2)
	sma := trend.NewSmaWithPeriod[T](b.Period)
	std := NewMovingStdWithPeriod[T](b.Period)

	middleBands := helper.Duplicate(
		sma.Compute(cs[0]),
		3,
	)

	std2s := helper.Duplicate(
		helper.MultiplyBy(
			std.Compute(cs[1]),
			2,
		),
		2,
	)

	upperBand := helper.Add(
		middleBands[0],
		std2s[0],
	)

	lowerBand := helper.Subtract(
		middleBands[1],
		std2s[1],
	)

	return upperBand, middleBands[2], lowerBand
}

// IdlePeriod is the initial period that Bollinger Bands won't yield any results.
func (b *BollingerBands[T]) IdlePeriod() int {
	return b.Period - 1
}
