// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import "github.com/cinar/indicator/v2/helper"

// Ad holds configuration parameters for calculating Accumulation/Distribution (A/D). It is a cumulative
// indicator that uses volume and price to assess whether an asset is being accumulated or distributed.
//
//	MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
//	MFV = MFM * Period Volume
//	AD = Previous AD + CMFV
//
// Example:
//
//	ad := volume.NewAd[float64]()
//	result := ad.Compute(highs, lows, closings, volumes)
type Ad[T helper.Number] struct {
	// Mfv is the MFV instance.
	Mfv *Mfv[T]
}

// NewAd function initializes a new A/D instance with the default parameters.
func NewAd[T helper.Number]() *Ad[T] {
	return &Ad[T]{
		Mfv: NewMfv[T](),
	}
}

// Compute function takes a channel of numbers and computes the A/D.
func (a *Ad[T]) Compute(highs, lows, closings, volumes <-chan T) <-chan T {
	//	MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
	//	MFV = MFM * Period Volume
	mfvs := a.Mfv.Compute(highs, lows, closings, volumes)

	//	AD = Previous AD + CMFV
	return helper.MapWithPrevious(mfvs, func(previous, current T) T {
		return previous + current
	}, 0)
}

// IdlePeriod is the initial period that A/D won't yield any results.
func (*Ad[T]) IdlePeriod() int {
	return 0
}
