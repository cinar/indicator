// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"
	"math"

	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultFisherPeriod is the default period for the Fisher Transform.
	DefaultFisherPeriod = 10

	// FisherClamp is the boundary value for clamping.
	FisherClamp = 0.999
)

// Fisher represents the configuration parameters for calculating the
// Fisher Transform. The Fisher Transform is a technical indicator
// that transforms prices into a normal distribution to identify
// price reversals.
//
//	x = 2 * ((close - min) / (max - min)) - 1
//	Fisher = 0.5 * ln((1 + x) / (1 - x))
//
// The clamped x value is bounded between -0.999 and +0.999 to prevent
// division by zero or logarithmic infinity errors.
//
// Example:
//
//	fisher := momentum.NewFisher[float64]()
//	result := fisher.Compute(closings)
type Fisher[T helper.Float] struct {
	// Period is the lookback period for min/max calculation.
	Period int

	// Max is the Moving Max instance.
	Max *trend.MovingMax[T]

	// Min is the Moving Min instance.
	Min *trend.MovingMin[T]
}

// NewFisher function initializes a new Fisher Transform instance.
func NewFisher[T helper.Float]() *Fisher[T] {
	return &Fisher[T]{
		Period: DefaultFisherPeriod,
		Max:    trend.NewMovingMaxWithPeriod[T](DefaultFisherPeriod),
		Min:    trend.NewMovingMinWithPeriod[T](DefaultFisherPeriod),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Fisher Transform.
func (f *Fisher[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	// Collect input to slice first to allow multiple independent channels
	values := helper.ChanToSlice(closings)

	// Create three independent channels from the slice
	input1 := helper.SliceToChanWithContext(ctx, values)
	input2 := helper.SliceToChanWithContext(ctx, values)
	input3 := helper.SliceToChanWithContext(ctx, values)

	// Compute min and max
	minValues := f.Min.ComputeWithContext(ctx, input1)
	maxValues := f.Max.ComputeWithContext(ctx, input2)

	// Align close values with min/max outputs
	alignedClosings := helper.SkipWithContext(ctx, input3, f.Period-1)

	// Compute: range = max - min
	rangeValues := helper.SubtractWithContext(ctx, maxValues, minValues)

	// Compute: close - min
	closeMinusMin := helper.SubtractWithContext(ctx, alignedClosings, minValues)

	// Compute: normalized = (close - min) / (max - min)
	normalized := helper.DivideWithContext(ctx, closeMinusMin, rangeValues)

	// Compute: x = 2 * normalized - 1
	x := helper.MapWithContext(ctx, normalized, func(v T) T {
		return 2*v - T(1)
	})

	// Clamp x to [-FisherClamp, FisherClamp] and compute Fisher
	result := helper.MapWithContext(ctx, x, func(v T) T {
		fx := float64(v)
		if fx > FisherClamp {
			fx = FisherClamp
		}
		if fx < -FisherClamp {
			fx = -FisherClamp
		}
		return T(0.5 * math.Log((1+fx)/(1-fx)))
	})

	return result
}

// IdlePeriod is the initial period that Fisher Transform won't yield any results.
func (f *Fisher[T]) IdlePeriod() int {
	// Min outputs after Period-1, Max outputs after Period-1
	// Close values need to skip Period-1
	// So total idle = Period-1 (from min/max) + Period-1 (from skip) = 2*Period-2
	return 2*f.Period - 2
}

// String is the string representation of the Fisher Transform.
func (f *Fisher[T]) String() string {
	return fmt.Sprintf("Fisher(%d)", f.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (f *Fisher[T]) Compute(closings <-chan T) <-chan T {
	return f.ComputeWithContext(context.Background(), closings)
}
