// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultCfoPeriod is the default CFO period of 14.
	DefaultCfoPeriod = 14
)

// Cfo represents the configuration parameters for calculating the
// Chande Forecast Oscillator (CFO). CFO is a momentum indicator
// that measures the difference between a security's price and its
// linear regression forecast.
//
//	CFO = ((Price - Forecast) / Price) * 100
//
// Example:
//
//	cfo := trend.NewCfo[float64]()
//	result := cfo.Compute(c)
type Cfo[T helper.Number] struct {
	// Mlr is the Moving Linear Regression instance.
	Mlr *Mlr[T]
}

// NewCfo function initializes a new CFO instance with the default parameters.
func NewCfo[T helper.Number]() *Cfo[T] {
	return NewCfoWithPeriod[T](DefaultCfoPeriod)
}

// NewCfoWithPeriod function initializes a new CFO instance with the given period.
func NewCfoWithPeriod[T helper.Number](period int) *Cfo[T] {
	return &Cfo[T]{
		Mlr: NewMlrWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the CFO.
func (c *Cfo[T]) Compute(closing <-chan T) <-chan T {
	closingSplices := helper.Duplicate(closing, 3)

	x := helper.Count(T(0), closingSplices[0])
	forecast := c.Mlr.Compute(x, closingSplices[1])

	closingPriceSplice := helper.Duplicate(helper.Skip(closingSplices[2], c.IdlePeriod()), 2)

	return helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(
				closingPriceSplice[0],
				forecast,
			),
			closingPriceSplice[1],
		),
		T(100),
	)
}

// IdlePeriod is the initial period that CFO won't yield any results.
func (c *Cfo[T]) IdlePeriod() int {
	return c.Mlr.IdlePeriod()
}
