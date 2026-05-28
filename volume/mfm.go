// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Mfm holds configuration parameters for calculating the Money Flow Multiplier (MFM),
// which adjusts volume based on the closing price's position within the high-low range:
//
//	MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
//
// - Positive MFM: Close in upper half of range, indicating buying pressure.
// - Negative MFM: Close in lower half of range, indicating selling pressure.
// - MFM of 1: Close equals high, strongest buying pressure.
// - MFM of -1: Close equals low, strongest selling pressure.
//
// Example:
//
//	mfm := volume.NewMfm[float64]()
//	result := mfm.Compute(highs, lows, closings)
type Mfm[T helper.Number] struct{}

// NewMfm function initializes a new MFM instance with the default parameters.
func NewMfm[T helper.Number]() *Mfm[T] {
	return &Mfm[T]{}
}

// ComputeWithContext function takes a channel of numbers and computes the MFM.
func (i *Mfm[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	highsSplice := helper.DuplicateWithContext(ctx, highs, 2)
	lowsSplice := helper.DuplicateWithContext(ctx, lows, 2)
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)

	return helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, helper.SubtractWithContext(ctx, closingsSplice[0], lowsSplice[0]),
		helper.SubtractWithContext(ctx, highsSplice[0], closingsSplice[1]),
	),
		helper.SubtractWithContext(ctx, highsSplice[1], lowsSplice[1]),
	)
}

// IdlePeriod is the initial period that MFM won't yield any results.
func (*Mfm[T]) IdlePeriod() int {
	return 0
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Mfm[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), highs, lows, closings)
}
