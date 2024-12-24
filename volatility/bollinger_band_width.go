// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
)

// BollingerBandWidth represents the configuration parameters for calculating the Bollinger Band Width.
// It measures the percentage difference between the upper band and the lower band. It decreases as
// Bollinger Bands narrows and increases as Bollinger Bands widens.
//
// During a period of rising price volatity the bandwidth widens, and during a period of low market
// volatity bandwidth contracts.
//
//	Band Width = (Upper Band - Lower Band) / Middle BollingerBandWidth
//
// Example:
//
//	bbw := NewBollingerBandWidth[float64]()
//	bbw.Compute(c)
type BollingerBandWidth[T helper.Number] struct {
	// Bollinger bands.
	BollingerBands *BollingerBands[T]
}

// NewBollingerBandWidth function initializes a new Bollinger Band Width instance with the default parameters.
func NewBollingerBandWidth[T helper.Number]() *BollingerBandWidth[T] {
	return &BollingerBandWidth[T]{
		BollingerBands: NewBollingerBands[T](),
	}
}

// Compute function takes a channel of numbers and computes the Bollinger Band Width.
func (b *BollingerBandWidth[T]) Compute(c <-chan T) <-chan T {
	upper, middle, lower := b.BollingerBands.Compute(c)

	return helper.Divide(
		helper.Subtract(upper, lower),
		middle,
	)
}

// IdlePeriod is the initial period that Bollinger Band Width won't yield any results.
func (b *BollingerBandWidth[T]) IdlePeriod() int {
	return b.BollingerBands.IdlePeriod()
}
