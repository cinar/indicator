// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultFiPeriod is the default period for the FI.
	DefaultFiPeriod = 13
)

// Fi holds configuration parameters for calculating the Force Index (FI). It uses the closing price and the volume to
// assess the power behind a move and identify turning points.
//
//	FI = EMA(period, (Current - Previous) * Volume)
//
// Example:
//
//	fi := volume.NewFi[float64]()
//	result := fi.Compute(closings, volumes)
type Fi[T helper.Float] struct {
	// Ema is the EMA instance.
	Ema *trend.Ema[T]
}

// NewFi function initializes a new FI instance with the default parameters.
func NewFi[T helper.Float]() *Fi[T] {
	return NewFiWithPeriod[T](DefaultFiPeriod)
}

// NewFiWithPeriod function initializes a new FI instance with the given period.
func NewFiWithPeriod[T helper.Float](period int) *Fi[T] {
	return &Fi[T]{
		Ema: trend.NewEmaWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the FI.
func (f *Fi[T]) ComputeWithContext(ctx context.Context, closings, volumes <-chan T) <-chan T {
	volumes = helper.SkipWithContext(ctx, volumes, 1)

	return f.Ema.ComputeWithContext(ctx, helper.MultiplyWithContext(ctx, helper.ChangeWithContext(ctx, closings, 1),
		volumes,
	),
	)
}

// IdlePeriod is the initial period that FI won't yield any results.
func (f *Fi[T]) IdlePeriod() int {
	return f.Ema.IdlePeriod() + 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (f *Fi[T]) Compute(closings, volumes <-chan T) <-chan T {
	return f.ComputeWithContext(context.Background(), closings, volumes)
}
