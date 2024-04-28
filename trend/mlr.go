// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

// Mlr represents the configuration parameters for calculating the Moving Linear Regression.
//
//	y = mx + b
//
// Example:
//
//	mlr := trend.NewMlrWithPeriod[float64](14)
//	rs := mlr.Compute(x , y)
type Mlr[T helper.Number] struct {
	// Mls is the Moving Least Square instance.
	Mls *Mls[T]
}

// NewMlrWithPeriod function initializes a new MLR instance with the given period.
func NewMlrWithPeriod[T helper.Number](period int) *Mlr[T] {
	return &Mlr[T]{
		Mls: NewMlsWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the MLR r.
func (m *Mlr[T]) Compute(x, y <-chan T) <-chan T {
	xSplice := helper.Duplicate(x, 2)

	ms, bs := m.Mls.Compute(xSplice[0], y)

	xSplice[1] = helper.Skip(xSplice[1], m.Mls.IdlePeriod())

	r := helper.Add(
		helper.Multiply(
			ms,
			xSplice[1],
		),
		bs,
	)

	return r
}

// IdlePeriod is the initial period that MLR won't yield any results.
func (m *Mlr[T]) IdlePeriod() int {
	return m.Mls.IdlePeriod()
}
