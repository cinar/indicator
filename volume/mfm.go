// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import "github.com/cinar/indicator/v2/helper"

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

// Compute function takes a channel of numbers and computes the MFM.
func (*Mfm[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)
	closingsSplice := helper.Duplicate(closings, 2)

	return helper.Divide(
		helper.Subtract(
			helper.Subtract(closingsSplice[0], lowsSplice[0]),
			helper.Subtract(highsSplice[0], closingsSplice[1]),
		),
		helper.Subtract(highsSplice[1], lowsSplice[1]),
	)
}

// IdlePeriod is the initial period that MFM won't yield any results.
func (*Mfm[T]) IdlePeriod() int {
	return 0
}
