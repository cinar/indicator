// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultCmfPeriod is the default period of CMF.
	DefaultCmfPeriod = 20
)

// Cmf holds configuration parameters for calculating the Chaikin Money Flow (CMF). It measures the amount
// of money flow volume over a given period.
//
//	MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
//	MFV = MFM * Volume
//	CMF = Sum(20, Money Flow Volume) / Sum(20, Volume)
//
// Example:
//
//	cmf := volume.NewCmf[float64]()
//	result := cmf.Compute(highs, lows, closings, volumes)
type Cmf[T helper.Number] struct {
	// Mfv is the MFV instance.
	Mfv *Mfv[T]

	// Sum is the Moving Sum instance.
	Sum *trend.MovingSum[T]
}

// NewCmf function initializes a new CMF instance with the default parameters.
func NewCmf[T helper.Number]() *Cmf[T] {
	return NewCmfWithPeriod[T](DefaultCmfPeriod)
}

// NewCmfWithPeriod function initializes a new CMF instance with the given period.
func NewCmfWithPeriod[T helper.Number](period int) *Cmf[T] {
	return &Cmf[T]{
		Mfv: NewMfv[T](),
		Sum: trend.NewMovingSumWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the CMF.
func (c *Cmf[T]) Compute(highs, lows, closings, volumes <-chan T) <-chan T {
	volumesSplice := helper.Duplicate(volumes, 2)

	//	MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
	//	MFV = MFM * Volume
	mfvs := c.Mfv.Compute(highs, lows, closings, volumesSplice[0])

	//	CMF = Sum(20, Money Flow Volume) / Sum(20, Volume)
	return helper.Divide(
		c.Sum.Compute(mfvs),
		c.Sum.Compute(volumesSplice[1]),
	)
}

// IdlePeriod is the initial period that MFV won't yield any results.
func (c *Cmf[T]) IdlePeriod() int {
	return c.Sum.IdlePeriod()
}
