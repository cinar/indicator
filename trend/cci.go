// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultCciPeriod is the default time period for CCI.
	DefaultCciPeriod = 20
)

// Cci represents the configuration parameters for calculating the Community Channel Index (CCI). CCI is a
// momentum-based oscillator used to help determine when an investment vehicle is reaching a condition of
// being overbought or oversold.
//
//	Moving Average = Sma(Period, Typical Price)
//	Mean Deviation = Sma(Period, Abs(Typical Price - Moving Average))
//	CCI = (Typical Price - Moving Average) / (0.015 * Mean Deviation)
//
// Example:
//
//	cmi := trend.NewCmi()
//	cmi.Period = 20
//	values = cmi.Compute(highs, lows, closings)
type Cci[T helper.Number] struct {
	// Time period.
	Period int
}

// NewCci function initializes a new CCI instance with the default parameters.
func NewCci[T helper.Number]() *Cci[T] {
	return &Cci[T]{
		Period: DefaultCciPeriod,
	}
}

// NewCciWithPeriod function initializes a new CCI instance with the given period.
func NewCciWithPeriod[T helper.Number](period int) *Cci[T] {
	return &Cci[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the CCI and the signal line.
func (c *Cci[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	typicalPrice := NewTypicalPrice[T]()
	sma1 := NewSmaWithPeriod[T](c.Period)
	sma2 := NewSmaWithPeriod[T](c.Period)

	tps := helper.Duplicate[T](
		typicalPrice.Compute(highs, lows, closings),
		3,
	)

	mas := helper.Duplicate[T](
		sma1.Compute(tps[0]),
		2,
	)

	tps[1] = helper.Skip(tps[1], sma1.Period-1)
	tps[2] = helper.Skip(tps[2], sma1.Period-1)

	md := sma2.Compute(
		helper.Abs(
			helper.Subtract(tps[1], mas[0]),
		),
	)

	mas[1] = helper.Skip(mas[1], sma2.Period-1)
	tps[2] = helper.Skip(tps[2], sma2.Period-1)

	multiplier := 0.015

	cci := helper.Divide(
		helper.Subtract(
			tps[2], mas[1],
		),
		helper.MultiplyBy[T](
			md,
			T(multiplier),
		),
	)

	return cci
}

// IdlePeriod is the initial period that CCI won't yield any results.
func (c *Cci[T]) IdlePeriod() int {
	return (c.Period * 2) - 2
}
