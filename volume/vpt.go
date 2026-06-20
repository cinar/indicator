// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"

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

// ComputeWithContext function takes a channel of numbers and computes the VPT.
func (i *Vpt[T]) ComputeWithContext(ctx context.Context, closings, volumes <-chan T) <-chan T {
	ratios := helper.MultiplyWithContext(ctx, helper.ChangeRatioWithContext(ctx, closings, 1),
		helper.SkipWithContext(ctx, volumes, 1),
	)

	return helper.MapWithPreviousWithContext(ctx, ratios, func(previous, current T) T {
		return previous + current
	}, 0)
}

// IdlePeriod is the initial period that VPT won't yield any results.
func (*Vpt[T]) IdlePeriod() int {
	return 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Vpt[T]) Compute(closings, volumes <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), closings, volumes)
}
