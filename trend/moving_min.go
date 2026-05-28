// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// MovingMin represents the configuration parameters for calculating the
// Moving Min over the specified period.
//
// Example:
type MovingMin[T helper.Number] struct {
	// Time period.
	Period int
}

// NewMovingMin function initializes a new Moving Min instance with the default parameters.
func NewMovingMin[T helper.Number]() *MovingMin[T] {
	return &MovingMin[T]{}
}

// NewMovingMinWithPeriod function initializes a new Moving Min instance with the given period.
func NewMovingMinWithPeriod[T helper.Number](period int) *MovingMin[T] {
	return &MovingMin[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the
// Moving Min over the specified period.
func (m *MovingMin[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	cs := helper.DuplicateWithContext(ctx, c, 2)
	cs[1] = helper.ShiftWithContext(ctx, cs[1], m.Period, 0)

	bst := helper.NewBst[T]()

	mins := helper.OperateWithContext(ctx, cs[0], cs[1], func(c, b T) T {
		bst.Insert(c)
		bst.Remove(b)
		return bst.Min()
	})

	return helper.SkipWithContext(ctx, mins, m.Period-1)
}

// IdlePeriod is the initial period that Mocing Min won't yield any results.
func (m *MovingMin[T]) IdlePeriod() int {
	return m.Period - 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *MovingMin[T]) Compute(c <-chan T) <-chan T {
	return m.ComputeWithContext(context.Background(), c)
}
