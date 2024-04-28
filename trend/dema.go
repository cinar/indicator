// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

// Dema represents the parameters for calculating the Double Exponential Moving Average (DEMA).
// A bullish cross occurs when DEMA with 5 days period moves above DEMA with 35 days period.
// A bearish cross occurs when DEMA with 35 days period moves above DEMA With 5 days period.
//
//	DEMA = (2 * EMA1(values)) - EMA2(EMA1(values))
//
// Example:
//
//	dema := trend.NewDema[float64]()
//	dema.Ema1.Period = 10
//	dema.Ema2.Period = 16
//
//	result := dema.Compute(input)
type Dema[T helper.Number] struct {
	// Ema1 represents the configuration parameters for
	// calculating the first EMA.
	Ema1 *Ema[T]

	// Ema2 represents the configuration parameters for
	// calculating the second EMA.
	Ema2 *Ema[T]
}

// NewDema function initializes a new DEMA instance
// with the default parameters.
func NewDema[T helper.Number]() *Dema[T] {
	return &Dema[T]{
		Ema1: NewEma[T](),
		Ema2: NewEma[T](),
	}
}

// Compute function takes a channel of numbers and computes the DEMA
// over the specified period.
func (d *Dema[T]) Compute(c <-chan T) <-chan T {
	ema1 := helper.Duplicate(d.Ema1.Compute(c), 2)
	ema2 := d.Ema2.Compute(ema1[1])

	doubleEma1 := helper.MultiplyBy(ema1[0], 2)
	doubleEma1 = helper.Buffered(doubleEma1, d.Ema2.Period)

	return helper.Subtract(doubleEma1, ema2)
}

// IdlePeriod is the initial period that DEMA won't yield any results.
func (d *Dema[T]) IdlePeriod() int {
	return d.Ema1.Period + d.Ema2.Period - 2
}
