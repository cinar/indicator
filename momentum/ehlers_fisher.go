// Copyright (c) 2021-2026 The Indicator Authors.
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
	// DefaultEhlersFisherPeriod is the default period for the Ehlers Fisher Transform.
	DefaultEhlersFisherPeriod = 10
)

// EhlersFisher represents the configuration parameters for calculating the
// canonical Fisher Transform, per John Ehlers' original Fisher Transform
// methodology ("Using the Fisher Transform", Stocks & Commodities magazine).
// Unlike [Fisher], it uses the median of high and low, and recursively
// smooths both the normalized price and the transformed output.
//
//	Price = (High + Low) / 2
//	Value1 = 0.33 * 2 * ((Price - MinL) / (MaxH - MinL) - 0.5) + 0.67 * Value1[previous]
//	Fisher = 0.5 * ln((1 + Value1) / (1 - Value1)) * 0.5 + 0.5 * Fisher[previous]
//
// Value1 is clamped to [-0.999, 0.999] before it is used in the logarithm,
// to prevent division by zero or logarithmic infinity errors. The very
// first emitted Value1 and Fisher, having no prior history, are computed
// with Value1[previous] and Fisher[previous] seeded at 0.
//
// Example:
//
//	ehlersFisher := momentum.NewEhlersFisher[float64]()
//	result := ehlersFisher.Compute(highs, lows)
type EhlersFisher[T helper.Float] struct {
	// Period is the lookback period for min/max calculation.
	Period int

	// Max is the Moving Max instance.
	Max *trend.MovingMax[T]

	// Min is the Moving Min instance.
	Min *trend.MovingMin[T]
}

// NewEhlersFisher function initializes a new Ehlers Fisher Transform instance
// with the default parameters.
func NewEhlersFisher[T helper.Float]() *EhlersFisher[T] {
	return NewEhlersFisherWithPeriod[T](DefaultEhlersFisherPeriod)
}

// NewEhlersFisherWithPeriod function initializes a new Ehlers Fisher
// Transform instance with the given period.
func NewEhlersFisherWithPeriod[T helper.Float](period int) *EhlersFisher[T] {
	return &EhlersFisher[T]{
		Period: period,
		Max:    trend.NewMovingMaxWithPeriod[T](period),
		Min:    trend.NewMovingMinWithPeriod[T](period),
	}
}

// ComputeWithContext function takes channels of high and low numbers and
// computes the Ehlers Fisher Transform.
func (e *EhlersFisher[T]) ComputeWithContext(ctx context.Context, highs, lows <-chan T) <-chan T {
	price := helper.DivideByWithContext(ctx, helper.AddWithContext(ctx, highs, lows), T(2))

	prices := helper.DuplicateWithContext(ctx, price, 3)
	price1, price2, price3 := prices[0], prices[1], prices[2]

	// minValues is used twice below (range and price-minus-min), so it
	// needs its own duplicated copy for the second use.
	minSplice := helper.DuplicateWithContext(ctx, e.Min.ComputeWithContext(ctx, price1), 2)
	maxValues := e.Max.ComputeWithContext(ctx, price2)

	// Align price values with min/max outputs.
	alignedPrices := helper.SkipWithContext(ctx, price3, e.Period-1)

	// Compute: range = max - min
	rangeValues := helper.SubtractWithContext(ctx, maxValues, minSplice[0])

	// Compute: price - min
	priceMinusMin := helper.SubtractWithContext(ctx, alignedPrices, minSplice[1])

	// Compute: normalized = (price - min) / (max - min), in [0, 1].
	normalized := helper.DivideWithContext(ctx, priceMinusMin, rangeValues)

	return e.smoothWithContext(ctx, normalized)
}

// smoothWithContext applies the recursive Value1 and Fisher smoothing
// described in the type's doc comment to the normalized [0, 1] price
// series, carrying the previous Value1 and Fisher across iterations.
func (e *EhlersFisher[T]) smoothWithContext(ctx context.Context, normalized <-chan T) <-chan T {
	result := make(chan T, cap(normalized))

	go func() {
		defer close(result)

		// Value1[previous] and Fisher[previous] are seeded at 0.
		var value1, fisher T

		for {
			select {
			case <-ctx.Done():
				return

			case n, ok := <-normalized:
				if !ok {
					return
				}

				v1 := 0.33*2*(float64(n)-0.5) + 0.67*float64(value1)

				if v1 > FisherClamp {
					v1 = FisherClamp
				} else if v1 < -FisherClamp {
					v1 = -FisherClamp
				}

				value1 = T(v1)
				fisher = T(0.5*math.Log((1+v1)/(1-v1))*0.5 + 0.5*float64(fisher))

				select {
				case <-ctx.Done():
					return
				case result <- fisher:
				}
			}
		}
	}()

	return result
}

// IdlePeriod is the initial period that Ehlers Fisher Transform won't yield any results.
func (e *EhlersFisher[T]) IdlePeriod() int {
	// Min, Max, and the aligned prices are each independently delayed
	// by Period-1 from the same input, not chained one after another, so
	// the delay doesn't compound.
	return e.Period - 1
}

// String is the string representation of the Ehlers Fisher Transform.
func (e *EhlersFisher[T]) String() string {
	return fmt.Sprintf("EhlersFisher(%d)", e.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (e *EhlersFisher[T]) Compute(highs, lows <-chan T) <-chan T {
	return e.ComputeWithContext(context.Background(), highs, lows)
}
