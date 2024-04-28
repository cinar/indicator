// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

// MovingMax represents the configuration parameters for calculating the
// Moving Max over the specified period.
//
// Example:
type MovingMax[T helper.Number] struct {
	// Time period.
	Period int
}

// NewMovingMax function initializes a new Moving Max instance with the default parameters.
func NewMovingMax[T helper.Number]() *MovingMax[T] {
	return &MovingMax[T]{}
}

// NewMovingMaxWithPeriod function initializes a new Moving Max instance with the given period.
func NewMovingMaxWithPeriod[T helper.Number](period int) *MovingMax[T] {
	max := NewMovingMax[T]()
	max.Period = period

	return max
}

// Compute function takes a channel of numbers and computes the
// Moving Max over the specified period.
func (m *MovingMax[T]) Compute(c <-chan T) <-chan T {
	cs := helper.Duplicate(c, 2)
	cs[1] = helper.Shift(cs[1], m.Period, 0)

	bst := helper.NewBst[T]()

	maxs := helper.Operate(cs[0], cs[1], func(c, b T) T {
		bst.Insert(c)
		bst.Remove(b)
		return bst.Max()
	})

	return helper.Skip(maxs, m.Period-1)
}

// IdlePeriod is the initial period that Mocing Max won't yield any results.
func (m *MovingMax[T]) IdlePeriod() int {
	return m.Period - 1
}
