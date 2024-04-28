// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

// Wma represents the configuration parameters for calculating the Weighted Moving Average (WMA).
// It calculates a moving average by putting more weight on recent data and less on past data.
//
//	WMA = ((Value1 * 1/N) + (Value2 * 2/N) + ...) / 2
type Wma[T helper.Number] struct {
	// Time period.
	Period int
}

// NewWmaWith function initializes a new WMA instance with the given parameters.
func NewWmaWith[T helper.Number](period int) *Wma[T] {
	return &Wma[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the WMA and the signal line.
func (w *Wma[T]) Compute(values <-chan T) <-chan T {
	window := helper.NewRing[T](w.Period)

	wmas := helper.Map(values, func(value T) T {
		window.Put(value)

		if !window.IsFull() {
			return T(0)
		}

		var sum T

		for i := 0; i < w.Period; i++ {
			sum += window.At(i) * T(i+1) / T(w.Period)
		}

		return sum / 2
	})

	wmas = helper.Skip(wmas, w.IdlePeriod())

	return wmas
}

// IdlePeriod is the initial period that WMA won't yield any results.
func (w *Wma[T]) IdlePeriod() int {
	return w.Period - 1
}

// String is the string representation of the WMA.
func (w *Wma[T]) String() string {
	return fmt.Sprintf("WMA(%d)", w.Period)
}
