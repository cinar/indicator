// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"math"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultAtrPeriod is the default period for the Average True Range (ATR).
	DefaultAtrPeriod = 14
)

// Atr represents the configuration parameters for calculating the Average True Range (ATR).
// It is a technical analysis indicator that measures market volatility by decomposing the
// entire range of stock prices for that period.
//
//	TR = Max((High - Low), (High - Previous Closing), (Previous Closing - Low))
//	ATR = MA TR
//
// By default, SMA is used as the MA.
//
// Example:
//
//	atr := volatility.NewAtr()
//	atr.Compute(highs, lows, closings)
type Atr[T helper.Number] struct {
	// Ma is the moving average for the ATR.
	Ma trend.Ma[T]
}

// NewAtr function initializes a new ATR instance with the default parameters.
func NewAtr[T helper.Number]() *Atr[T] {
	return NewAtrWithPeriod[T](DefaultAtrPeriod)
}

// NewAtrWithPeriod function initializes a new ATR instance with the given period.
func NewAtrWithPeriod[T helper.Number](period int) *Atr[T] {
	return NewAtrWithMa(trend.NewSmaWithPeriod[T](period))
}

// NewAtrWithMa function initializes a new ATR instance with the given moving average instance.
func NewAtrWithMa[T helper.Number](ma trend.Ma[T]) *Atr[T] {
	return &Atr[T]{
		Ma: ma,
	}
}

// Compute function takes a channel of numbers and computes the ATR over the specified period.
func (a *Atr[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	// Use previous closing by skipping highs and lows by one.
	highs = helper.Skip(highs, 1)
	lows = helper.Skip(lows, 1)

	tr := helper.Operate3(highs, lows, closings, func(high, low, closing T) T {
		return T(math.Max(float64(high-low), math.Max(float64(high-closing), float64(closing-low))))
	})

	atr := a.Ma.Compute(tr)

	return atr
}

// IdlePeriod is the initial period that Acceleration Bands won't yield any results.
func (a *Atr[T]) IdlePeriod() int {
	// Ma idle period and for using the previous closing.
	return a.Ma.IdlePeriod() + 1
}
