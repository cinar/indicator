// Copyright (c) 2021-2026 Onur Cinar.
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
type Obv[T helper.Number] struct{}

// NewObv function initializes a new OBV instance with the default parameters.
func NewObv[T helper.Number]() *Obv[T] {
	return &Obv[T]{}
}

// ComputeWithContext function takes a channel of numbers and computes the OBV.
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *Obv[T]) Compute(closings, volumes <-chan T) <-chan T {
	return i.ComputeWithContext(context.Background(), closings, volumes)
}
