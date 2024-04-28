// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import "github.com/cinar/indicator/v2/helper"

// Mfv holds configuration parameters for calculating Money Flow Volume (MFV), a volume-based indicator that
// incorporates the Money Flow Multiplier (MFM) to gauge the intensity of buying and selling pressure. MFV
// reflects the cumulative volume adjusted by MFM, with higher values indicating stronger buying pressure
// and lower values suggesting selling dominance. MFV highlights periods of significant volume-driven
// price action, offering insights into potential trend strength and reversals.
//
//	MFV = MFM * Volume
//
// Example:
//
//	mfv := volume.NewMfv[float64]()
//	result := mfv.Compute(highs, lows, closings, volumes)
type Mfv[T helper.Number] struct {
	// Mfm is the MFM instance.
	Mfm *Mfm[T]
}

// NewMfv function initializes a new MFV instance with the default parameters.
func NewMfv[T helper.Number]() *Mfv[T] {
	return &Mfv[T]{
		Mfm: NewMfm[T](),
	}
}

// Compute function takes a channel of numbers and computes the MFV.
func (m *Mfv[T]) Compute(highs, lows, closings, volumes <-chan T) <-chan T {
	return helper.Multiply(
		m.Mfm.Compute(highs, lows, closings),
		volumes,
	)
}

// IdlePeriod is the initial period that MFV won't yield any results.
func (*Mfv[T]) IdlePeriod() int {
	return 0
}
