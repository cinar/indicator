// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"math"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

// TrueRange represents the parameters for calculating the True Range (TR).
// It evaluates the greatest distance covered by price in a single period,
// accounting for gaps.
//
//	TR = Max((High - Low), (High - Previous Closing), (Previous Closing - Low))
//
// Example:
//
//	tr := volatility.NewTrueRange[float64]()
//	result := tr.Compute(highs, lows, closings)
type TrueRange[T helper.Number] struct{}

// NewTrueRange function initializes a new TrueRange instance.
func NewTrueRange[T helper.Number]() *TrueRange[T] {
	return &TrueRange[T]{}
}

// ComputeWithContext function takes channels of highs, lows, and closings and computes the True Range.
func (tr *TrueRange[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	highs = helper.SkipWithContext(ctx, highs, 1)
	lows = helper.SkipWithContext(ctx, lows, 1)

	return helper.Operate3WithContext(ctx, highs, lows, closings, func(high, low, closing T) T {
		return T(math.Max(float64(high-low), math.Max(float64(high-closing), float64(closing-low))))
	})
}

// IdlePeriod is the initial period that TrueRange won't yield any results.
func (tr *TrueRange[T]) IdlePeriod() int {
	return 1
}

// String is the string representation of the TrueRange.
func (tr *TrueRange[T]) String() string {
	return "TR"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (tr *TrueRange[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return tr.ComputeWithContext(context.Background(), highs, lows, closings)
}
