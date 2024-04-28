// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

// Mls represents the configuration parameters for calculating the Moving Least Square (MLS). It is a regression
// analysis to determine the line of best fit for the given set of data.
//
//	y = mx + b
//	b = y-intercept
//	y = slope
//
//	m = (period * sumXY - sumX * sumY) / (period * sumX2 - sumX * sumX)
//	b = (sumY - m * sumX) / period
//
// Example:
//
//	mls := trend.NewMlsWithPeriod[float64](14)
//	ms, bs := mls.Compute(x , y)
type Mls[T helper.Number] struct {
	// Sum is the moving sum instance.
	Sum *MovingSum[T]
}

// NewMlsWithPeriod function initializes a new MLS instance with the given period.
func NewMlsWithPeriod[T helper.Number](period int) *Mls[T] {
	return &Mls[T]{
		Sum: NewMovingSumWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the MLS m and b.
func (m *Mls[T]) Compute(x, y <-chan T) (<-chan T, <-chan T) {
	xSplice := helper.Duplicate(x, 3)
	ySplice := helper.Duplicate(y, 2)

	sumXY := m.Sum.Compute(
		helper.Operate(xSplice[0], ySplice[0], func(a, b T) T {
			return a * b
		}),
	)

	sumXSplice := helper.Duplicate(
		m.Sum.Compute(xSplice[1]),
		4,
	)

	sumYSplice := helper.Duplicate(
		m.Sum.Compute(ySplice[1]),
		2,
	)

	sumX2 := m.Sum.Compute(
		helper.Pow(xSplice[2], 2),
	)

	// m = (period * sumXY - sumX * sumY) / (period * sumX2 - sumX * sumX)
	mSplice := helper.Duplicate(
		helper.Divide(
			helper.Subtract(
				helper.MultiplyBy(
					sumXY,
					T(m.Sum.Period),
				),
				helper.Multiply(
					sumXSplice[0],
					sumYSplice[0],
				),
			),
			helper.Subtract(
				helper.MultiplyBy(
					sumX2,
					T(m.Sum.Period),
				),
				helper.Multiply(
					sumXSplice[1],
					sumXSplice[2],
				),
			),
		),
		2,
	)

	// b = (sumY - m * sumX) / period
	b := helper.DivideBy(
		helper.Subtract(
			sumYSplice[1],
			helper.Multiply(
				mSplice[1],
				sumXSplice[3],
			),
		),
		T(m.Sum.Period),
	)

	return mSplice[0], b
}

// IdlePeriod is the initial period that MLS won't yield any results.
func (m *Mls[T]) IdlePeriod() int {
	return m.Sum.IdlePeriod()
}
