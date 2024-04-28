// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

// MovingSum represents the configuration parameters for calculating the Moving Sum over the specified period.
//
// Example:
//
//	sum := trend.NewMovingSum[float64]()
//	sum.Period = 20
type MovingSum[T helper.Number] struct {
	// Time period.
	Period int
}

// NewMovingSum function initializes a new Moving Sum instance with the default parameters.
func NewMovingSum[T helper.Number]() *MovingSum[T] {
	return NewMovingSumWithPeriod[T](1)
}

// NewMovingSumWithPeriod function initializes a new Moving Sum instance with the given period.
func NewMovingSumWithPeriod[T helper.Number](period int) *MovingSum[T] {
	return &MovingSum[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the
// Moving Sum over the specified period.
func (m *MovingSum[T]) Compute(c <-chan T) <-chan T {
	cs := helper.Duplicate(c, 2)
	cs[1] = helper.Shift(cs[1], m.Period, 0)

	sum := T(0)

	sums := helper.Operate(cs[0], cs[1], func(c, b T) T {
		sum = sum + c - b
		return sum
	})

	return helper.Skip(sums, m.Period-1)
}

// IdlePeriod is the initial period that Moving Sum won't yield any results.
func (m *MovingSum[T]) IdlePeriod() int {
	return m.Period - 1
}
