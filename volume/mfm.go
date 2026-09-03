// Copyright (c) 2021-2026 The Indicator Authors.
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
// On a flat bar (High == Low), the close's position within the range is undefined
// (0/0). MFM returns 0, the neutral point of its [-1, 1] scale, matching the
// InternalBarStrength precedent for the identical High-Low denominator. Since MFM
// feeds Mfv, Ad, and Cmf, a neutral 0 flows through as "zero contribution" to all
// of them.
//
// Example:
//
//	mfm := volume.NewMfm[float64]()
//	result := mfm.Compute(highs, lows, closings)
type Mfm[T helper.Float] struct{}

// NewMfm function initializes a new MFM instance with the default parameters.
func NewMfm[T helper.Float]() *Mfm[T] {
	return &Mfm[T]{}
}

// ComputeWithContext function takes a channel of numbers and computes the MFM.
func (i *Mfm[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	return helper.Operate3WithContext(ctx, highs, lows, closings, func(high, low, closing T) T {
		denom := high - low
		if denom == 0 {
			return 0
		}

		return ((closing - low) - (high - closing)) / denom
	})
}

// IdlePeriod is the initial period that MFM won't yield any results.
func (*Mfm[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the MFM.
func (*Mfm[T]) String() string {
	return "MFM"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Mfm[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), highs, lows, closings)
}
