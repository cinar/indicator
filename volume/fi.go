// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
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
type Fi[T helper.Number] struct {
	// Ema is the EMA instance.
	Ema *trend.Ema[T]
}

// NewFi function initializes a new FI instance with the default parameters.
func NewFi[T helper.Number]() *Fi[T] {
	return NewFiWithPeriod[T](DefaultFiPeriod)
}

// NewFiWithPeriod function initializes a new FI instance with the given period.
func NewFiWithPeriod[T helper.Number](period int) *Fi[T] {
	return &Fi[T]{
		Ema: trend.NewEmaWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the FI.
func (f *Fi[T]) Compute(closings, volumes <-chan T) <-chan T {
	return f.Ema.Compute(
		helper.Multiply(
			helper.Change(closings, 1),
			volumes,
		),
	)
}

// IdlePeriod is the initial period that FI won't yield any results.
func (f *Fi[T]) IdlePeriod() int {
	return f.Ema.IdlePeriod() + 1
}
