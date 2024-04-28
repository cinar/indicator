// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

// Tema represents the configuration parameters for calculating the
// Triple Exponential Moving Average (TEMA).
//
//	TEMA = (3 * EMA1) - (3 * EMA2) + EMA3
//	EMA1 = EMA(values)
//	EMA2 = EMA(EMA1)
//	EMA3 = EMA(EMA2)
type Tema[T helper.Number] struct {
	Ema1 *Ema[T]
	Ema2 *Ema[T]
	Ema3 *Ema[T]
}

// NewTema function initializes a new TEMA instance
// with the default parameters.
func NewTema[T helper.Number]() *Tema[T] {
	return &Tema[T]{
		Ema1: NewEma[T](),
		Ema2: NewEma[T](),
		Ema3: NewEma[T](),
	}
}

// Compute function takes a channel of numbers and computes the TEMA
// and the signal line.
func (t *Tema[T]) Compute(c <-chan T) <-chan T {
	ema1 := helper.Duplicate(
		t.Ema1.Compute(c),
		2,
	)

	ema2 := helper.Duplicate(
		t.Ema2.Compute(ema1[0]),
		2,
	)

	ema1[1] = helper.Skip(ema1[1], t.Ema2.Period-1)

	ema3 := t.Ema3.Compute(ema2[0])
	ema1[1] = helper.Skip(ema1[1], t.Ema3.Period-1)
	ema2[1] = helper.Skip(ema2[1], t.Ema3.Period-1)

	tema := helper.Add(
		helper.Subtract(
			helper.MultiplyBy(ema1[1], 3),
			helper.MultiplyBy(ema2[1], 3),
		),
		ema3,
	)

	return tema
}

// IdlePeriod is the initial period that TEMA won't yield any results.
func (t *Tema[T]) IdlePeriod() int {
	return t.Ema1.Period + t.Ema2.Period + t.Ema3.Period - 3
}
