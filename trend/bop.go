// Copyright (c) 2021-2026 Onur Cinar.
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
type Bop[T helper.Number] struct{}

// NewBop function initializes a new BOP instance
// with the default parameters.
func NewBop[T helper.Number]() *Bop[T] {
	return &Bop[T]{}
}

// ComputeWithContext processes a channel of open, high, low, and close values,
// computing the BOP for each entry.
func (i *Bop[T]) ComputeWithContext(ctx context.Context, opening, high, low, closing <-chan T) <-chan T {
	return helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, closing, opening), helper.SubtractWithContext(ctx, high, low))
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Bop[T]) Compute(opening, high, low, closing <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), opening, high, low, closing)
}
