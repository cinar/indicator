// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// InternalBarStrength represents the parameters for calculating the
// Internal Bar Strength (IBS). It tracks price location within a daily
// high-low range.
//
//	IBS = (Close - Low) / (High - Low)
//
// Example:
//
//	ibs := momentum.NewInternalBarStrength[float64]()
//	result := ibs.Compute(highs, lows, closings)
type InternalBarStrength[T helper.Number] struct{}

// NewInternalBarStrength function initializes a new InternalBarStrength instance.
func NewInternalBarStrength[T helper.Number]() *InternalBarStrength[T] {
	return &InternalBarStrength[T]{}
}

// ComputeWithContext function takes channels of highs, lows, and closings and computes the IBS.
func (ibs *InternalBarStrength[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	return helper.Operate3WithContext(ctx, highs, lows, closings, func(high, low, closing T) T {
		denom := high - low
		if denom == 0 {
			return 0
		}
		return (closing - low) / denom
	})
}

// IdlePeriod is the initial period that InternalBarStrength won't yield any results.
func (ibs *InternalBarStrength[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the InternalBarStrength.
func (ibs *InternalBarStrength[T]) String() string {
	return "IBS"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (ibs *InternalBarStrength[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return ibs.ComputeWithContext(context.Background(), highs, lows, closings)
}
