// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"
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
	wma1 *Wma[T]

	// Second WMA.
	wma2 *Wma[T]

	// Third WMA.
	wma3 *Wma[T]
}

// NewHmaWithPeriod function initializes a new HMA instance with the given parameters.
func NewHmaWithPeriod[T helper.Number](period int) *Hma[T] {
	return &Hma[T]{
		wma1: NewWmaWith[T](int(math.Round(float64(period) / 2))),
		wma2: NewWmaWith[T](period),
		wma3: NewWmaWith[T](int(math.Round(math.Sqrt(float64(period))))),
	}
}

// Compute function takes a channel of numbers and computes the HMA and the signal line.
func (h *Hma[T]) Compute(values <-chan T) <-chan T {
	valuesSplice := helper.Duplicate(values, 2)

	//	WMA1 = WMA(period/2 , values)
	wmas1 := h.wma1.Compute(valuesSplice[0])

	//	WMA2 = WMA(period, values)
	wmas2 := h.wma2.Compute(valuesSplice[1])

	wmas1 = helper.Skip(wmas1, h.wma2.IdlePeriod()-h.wma1.IdlePeriod())

	// WMA3 = WMA(sqrt(period), (2 * WMA1) - WMA2)
	wmas3 := h.wma3.Compute(
		helper.Subtract(
			helper.MultiplyBy(
				wmas1,
				2,
			),
			wmas2,
		),
	)

	// HMA = WMA3
	return wmas3
}

// IdlePeriod is the initial period that HMA won't yield any results.
func (h *Hma[T]) IdlePeriod() int {
	return h.wma2.IdlePeriod() + h.wma3.IdlePeriod()
}

// String is the string representation of the HMA.
func (h *Hma[T]) String() string {
	return fmt.Sprintf("HMA(%d)", h.wma2.Period)
}
