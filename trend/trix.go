// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultTrixPeriod is the default time period for TRIX.
	DefaultTrixPeriod = 15
)

// Trix represents the configuration parameters for calculating the Triple Exponential Average (TRIX).
// TRIX indicator is an oscillator used to identify oversold and overbought markets, and it can also
// be used as a momentum indicator. Like many oscillators, TRIX oscillates around a zero line.
//
//	EMA1 = EMA(period, values)
//	EMA2 = EMA(period, EMA1)
//	EMA3 = EMA(period, EMA2)
//	TRIX = (EMA3 - Previous EMA3) / Previous EMA3
//
// Example:
//
//	trix := trend.NewTrix[float64]()
//	result := trix.Compute(values)
type Trix[T helper.Number] struct {
	// Time period.
	Period int
}

// NewTrix function initializes a new TRIX instance with the default parameters.
func NewTrix[T helper.Number]() *Trix[T] {
	return &Trix[T]{
		Period: DefaultTrixPeriod,
	}
}

// Compute function takes a channel of numbers and computes the TRIX and the signal line.
func (t *Trix[T]) Compute(c <-chan T) <-chan T {
	ema1 := NewEmaWithPeriod[T](t.Period)
	ema2 := NewEmaWithPeriod[T](t.Period)
	ema3 := NewEmaWithPeriod[T](t.Period)

	emas := ema3.Compute(
		ema2.Compute(
			ema1.Compute(c),
		),
	)

	trix := helper.ChangeRatio[T](emas, 1)

	return trix
}

// IdlePeriod is the initial period that TRIX won't yield any results.
func (t *Trix[T]) IdlePeriod() int {
	return (t.Period * 3) - 3 + 1
}
