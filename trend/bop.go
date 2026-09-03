// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Bop gauges the strength of buying and selling forces using
// the Balance of Power (BoP) indicator. A positive BoP value
// suggests an upward trend, while a negative value indicates
// a downward trend. A BoP value of zero implies equilibrium
// between the two forces.
//
//	Formula: BOP = (Closing - Opening) / (High - Low)
//
// A zero-range bar (High == Low) also forces Open == Close == High == Low,
// so the numerator is zero too; BOP is defined as 0 (equilibrium, its
// natural centered value) instead of propagating the resulting 0/0 NaN.
type Bop[T helper.Float] struct{}

// NewBop function initializes a new BOP instance
// with the default parameters.
func NewBop[T helper.Float]() *Bop[T] {
	return &Bop[T]{}
}

// ComputeWithContext processes a channel of open, high, low, and close values,
// computing the BOP for each entry.
func (i *Bop[T]) ComputeWithContext(ctx context.Context, opening, high, low, closing <-chan T) <-chan T {
	numerator := helper.SubtractWithContext(ctx, closing, opening)
	denominator := helper.SubtractWithContext(ctx, high, low)

	return helper.OperateWithContext(ctx, numerator, denominator, func(num, denom T) T {
		if denom == 0 {
			return 0
		}

		return num / denom
	})
}

// IdlePeriod is the initial period that BOP won't yield any results.
func (*Bop[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the BOP.
func (*Bop[T]) String() string {
	return "BOP"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Bop[T]) Compute(opening, high, low, closing <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), opening, high, low, closing)
}
