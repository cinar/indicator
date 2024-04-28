// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultVwapPeriod is the default period for the VWAP.
	DefaultVwapPeriod = 14
)

// Vwap holds configuration parameters for calculating the Volume Weighted Average Price (VWAP). It provides the
// average price the asset has traded.
//
//	VWAP = Sum(Closing * Volume) / Sum(Volume)
//
// Example:
//
//	vwap := volume.NewVwap[float64]()
//	result := vwap.Compute(closings, volumes)
type Vwap[T helper.Number] struct {
	// Sum is the Moving Sum instance.
	Sum *trend.MovingSum[T]
}

// NewVwap function initializes a new VWAP instance with the default parameters.
func NewVwap[T helper.Number]() *Vwap[T] {
	return NewVwapWithPeriod[T](DefaultVwapPeriod)
}

// NewVwapWithPeriod function initializes a new VWAP instance with the given period.
func NewVwapWithPeriod[T helper.Number](period int) *Vwap[T] {
	return &Vwap[T]{
		Sum: trend.NewMovingSumWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the VWAP.
func (v *Vwap[T]) Compute(closings, volumes <-chan T) <-chan T {
	volumesSplice := helper.Duplicate(volumes, 2)

	return helper.Divide(
		v.Sum.Compute(
			helper.Multiply(
				closings,
				volumesSplice[0],
			),
		),
		v.Sum.Compute(
			volumesSplice[1],
		),
	)
}

// IdlePeriod is the initial period that VWAP won't yield any results.
func (v *Vwap[T]) IdlePeriod() int {
	return v.Sum.IdlePeriod()
}
