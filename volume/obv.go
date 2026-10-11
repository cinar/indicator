// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Obv holds configuration parameters for calculating the On-Balance Volume (OBV). It is a technical trading momentum
// indicator that uses volume flow to predict changes in asset price.
//
//	Foreach Closing:
//		If Closing[i] > Closing[i-1], OBV[i] = OBV[i-1] + Volume[i]
//		If Closing[i] = Closing[i-1], OBV[i] = OBV[i-1]
//		If Closing[i] < Closing[i-1], OBV[i] = OBV[i-1] - Volume[i]
//
// Example:
//
//	obv := volume.NewObv[float64]()
//	result := obv.Compute(closings, volumes)
//
// Note that the first emitted value, OBV[0], includes the full first-bar volume rather than starting at exactly
// 0. This is because there is no true "previous close" for the very first bar, so the zero-valued previousClosing
// causes the first comparison to read as an increase by default. This is inconsequential for typical OBV usage,
// such as slope or trend analysis, since it only introduces a constant offset to the series rather than changing
// its shape.
type Obv[T helper.Number] struct{}

// NewObv function initializes a new OBV instance with the default parameters.
func NewObv[T helper.Number]() *Obv[T] {
	return &Obv[T]{}
}

// ComputeWithContext function takes a channel of numbers and computes the OBV.
//
// Note that the first result includes the full first-bar volume rather than 0, since previousClosing starts at
// its zero value and there is no real prior close to compare against for the first bar.
func (i *Obv[T]) ComputeWithContext(ctx context.Context, closings, volumes <-chan T) <-chan T {
	var previousClosing T
	var previousObv T

	return helper.OperateWithContext(ctx, closings, volumes, func(closing, volume T) T {
		currentObv := previousObv

		if closing > previousClosing {
			currentObv += volume
		} else if closing < previousClosing {
			currentObv -= volume
		}

		previousClosing = closing
		previousObv = currentObv
		return currentObv
	})
}

// IdlePeriod is the initial period that OBV won't yield any results.
func (*Obv[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the OBV.
func (*Obv[T]) String() string {
	return "OBV"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Obv[T]) Compute(closings, volumes <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), closings, volumes)
}
