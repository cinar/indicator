// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package trend

import (
	"math"

	"github.com/cinar/indicator/v2/helper"
)

// Hma represents the configuration parameters for calculating the Hull Moving Average (HMA). Developed by
// Alan Hull in 2005, HMA attempts to minimize the lag of a traditional moving average.
//
//	WMA1 = WMA(period/2 , values)
//	WMA2 = WMA(period, values)
//	WMA3 = WMA(sqrt(period), (2 * WMA1) - WMA2)
//	HMA = WMA3
type Hma[T helper.Number] struct {
	// First WMA.
	Wma1 *Wma[T]

	// Second WMA.
	Wma2 *Wma[T]

	// Third WMA.
	Wma3 *Wma[T]
}

// NewHmaWith function initializes a new HMA instance with the given parameters.
func NewHmaWith[T helper.Number](period int) *Hma[T] {
	return &Hma[T]{
		Wma1: NewWmaWith[T](int(math.Round(float64(period) / 2))),
		Wma2: NewWmaWith[T](period),
		Wma3: NewWmaWith[T](int(math.Round(math.Sqrt(float64(period))))),
	}
}

// Compute function takes a channel of numbers and computes the HMA and the signal line.
func (h *Hma[T]) Compute(values <-chan T) <-chan T {
	valuesSplice := helper.Duplicate(values, 2)

	//	WMA1 = WMA(period/2 , values)
	wma1 := h.Wma1.Compute(valuesSplice[0])

	//	WMA2 = WMA(period, values)
	wma2 := h.Wma2.Compute(valuesSplice[1])

	wma1 = helper.Skip(wma1, h.Wma2.IdlePeriod()-h.Wma1.IdlePeriod())

	// WMA3 = WMA(sqrt(period), (2 * WMA1) - WMA2)
	wma3 := h.Wma3.Compute(
		helper.Subtract(
			helper.MultiplyBy(
				wma1,
				2,
			),
			wma2,
		),
	)

	// HMA = WMA3
	return wma3
}

// IdlePeriod is the initial period that HMA won't yield any results.
func (h *Hma[T]) IdlePeriod() int {
	return h.Wma2.IdlePeriod() + h.Wma3.IdlePeriod()
}
