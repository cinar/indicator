// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

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

// ComputeWithContext function takes a channel of numbers and computes the
// Moving Sum over the specified period.
//
// The running sum is accumulated with Neumaier (improved Kahan) compensated
// summation so that floating-point rounding error stays bounded regardless
// of how long the input series is, rather than compounding on every
// add/subtract step. Neumaier's variant, unlike classic Kahan, stays
// accurate even when an added term is larger in magnitude than the running
// sum, which happens whenever the window sum is near zero — e.g. a Cmf or
// Mfi moving sum of signed money flow in a quiet market. For integer T the
// compensation term is always zero, so behavior is unchanged.
func (m *MovingSum[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	cs := helper.DuplicateWithContext(ctx, c, 2)
	cs[1] = helper.ShiftWithContext(ctx, cs[1], m.Period, 0)

	sum := T(0)
	comp := T(0)

	sums := helper.OperateWithContext(ctx, cs[0], cs[1], func(c, b T) T {
		sum, comp = neumaierAdd(sum, comp, c)
		sum, comp = neumaierAdd(sum, comp, -b)

		return sum + comp
	})

	return helper.SkipWithContext(ctx, sums, m.Period-1)
}

// neumaierAdd adds value to sum using Neumaier compensated summation,
// returning the updated running sum and compensation term. The
// numerically corrected total is sum+comp.
func neumaierAdd[T helper.Number](sum, comp, value T) (T, T) {
	t := sum + value

	if absT(sum) >= absT(value) {
		comp += (sum - t) + value
	} else {
		comp += (value - t) + sum
	}

	return t, comp
}

// absT returns the absolute value of v.
func absT[T helper.Number](v T) T {
	if v < 0 {
		return -v
	}

	return v
}

// IdlePeriod is the initial period that Moving Sum won't yield any results.
func (m *MovingSum[T]) IdlePeriod() int {
	return m.Period - 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *MovingSum[T]) Compute(c <-chan T) <-chan T {
	return m.ComputeWithContext(context.Background(), c)
}
