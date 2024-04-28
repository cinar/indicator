// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
)

// Vpt holds configuration parameters for calculating the Volume Price Trend (VPT). It provides a correlation
// between the volume and the price.
//
//	VPT = Previous VPT + (Volume * (Current Closing - Previous Closing) / Previous Closing)
//
// Example:
//
//	vpt := volume.NewVpt[float64]()
//	result := vpt.Compute(closings, volumes)
type Vpt[T helper.Number] struct{}

// NewVpt function initializes a new VPT instance with the default parameters.
func NewVpt[T helper.Number]() *Vpt[T] {
	return &Vpt[T]{}
}

// Compute function takes a channel of numbers and computes the VPT.
func (*Vpt[T]) Compute(closings, volumes <-chan T) <-chan T {
	ratios := helper.Multiply(
		helper.ChangeRatio(closings, 1),
		helper.Skip(volumes, 1),
	)

	return helper.MapWithPrevious(ratios, func(previous, current T) T {
		return previous + current
	}, 0)
}

// IdlePeriod is the initial period that VPT won't yield any results.
func (*Vpt[T]) IdlePeriod() int {
	return 1
}
