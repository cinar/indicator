// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"context"

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

// ComputeWithContext function takes a channel of numbers and computes the Keltner Channel over the specified period.
func (k *KeltnerChannel[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) (<-chan T, <-chan T, <-chan T) {
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)

	//	2 * ATR(period, highs, lows, closings)
	atrs := helper.DuplicateWithContext(ctx, helper.MultiplyByWithContext(ctx, k.Atr.ComputeWithContext(ctx, highs, lows, closingsSplice[0]),
		2,
	),
		2,
	)

	//	Middle Line = EMA(period, closings)
	middles := helper.DuplicateWithContext(ctx, helper.SkipWithContext(ctx, k.Ema.ComputeWithContext(ctx, closingsSplice[1]),
		k.Atr.IdlePeriod()-k.Ema.IdlePeriod(),
	),
		3,
	)

	//	Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
	upper := helper.AddWithContext(ctx, middles[0], atrs[0])

	//	Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
	lower := helper.SubtractWithContext(ctx, middles[1], atrs[1])

	return upper, middles[2], lower
}

// IdlePeriod is the initial period that Keltner Channel won't yield any results.
func (k *KeltnerChannel[T]) IdlePeriod() int {
	return k.Atr.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (k *KeltnerChannel[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T, <-chan T) {
	return k.ComputeWithContext(context.Background(), highs, lows, closings)
}
