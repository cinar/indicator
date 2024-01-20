// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultKeltnerChannelPeriod is the default period for the Keltner Channel.
	DefaultKeltnerChannelPeriod = 20
)

// KeltnerChannel represents the configuration parameters for calculating the Keltner Channel (KC). It provides
// volatility-based bands that are placed on either side of an asset's price and can aid in determining the
// direction of a trend.
//
//	Middle Line = EMA(period, closings)
//	Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
//	Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
//
// Example:
//
//	dc := volatility.NewKeltnerChannel[float64]()
//	result := dc.Compute(highs, lows, closings)
type KeltnerChannel[T helper.Number] struct {
	// Atr is the ATR instance.
	Atr *Atr[T]

	// Ema is the EMA instance.
	Ema *trend.Ema[T]
}

// NewKeltnerChannel function initializes a new Keltner Channel instance with the default parameters.
func NewKeltnerChannel[T helper.Number]() *KeltnerChannel[T] {
	return NewKeltnerChannelWithPeriod[T](DefaultKeltnerChannelPeriod)
}

// NewKeltnerChannelWithPeriod function initializes a new Keltner Channel instance with the given period.
func NewKeltnerChannelWithPeriod[T helper.Number](period int) *KeltnerChannel[T] {
	return &KeltnerChannel[T]{
		Atr: NewAtrWithPeriod[T](period),
		Ema: trend.NewEmaWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the Keltner Channel over the specified period.
func (k *KeltnerChannel[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T, <-chan T) {
	closingsSplice := helper.Duplicate(closings, 2)

	//	2 * ATR(period, highs, lows, closings)
	atrs := helper.Duplicate(
		helper.MultiplyBy(
			k.Atr.Compute(highs, lows, closingsSplice[0]),
			2,
		),
		2,
	)

	//	Middle Line = EMA(period, closings)
	middles := helper.Duplicate(
		k.Ema.Compute(closingsSplice[1]),
		3,
	)

	//	Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
	upper := helper.Add(middles[0], atrs[0])

	//	Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
	lower := helper.Subtract(middles[1], atrs[1])

	return upper, middles[2], lower
}

// IdlePeriod is the initial period that Keltner Channel won't yield any results.
func (k *KeltnerChannel[T]) IdlePeriod() int {
	return k.Atr.IdlePeriod()
}
