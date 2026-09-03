// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultTsiFirstSmoothingPeriod is the default first smoothing period of 25.
	DefaultTsiFirstSmoothingPeriod = 25

	// DefaultTsiSecondSmoothingPeriod is the default second smoothing period of 13.
	DefaultTsiSecondSmoothingPeriod = 13
)

// Tsi represents the parameters needed to calculate the True Strength Index (TSI). It is a technical momentum
// oscillator used in financial analysis. The TSI helps identify trends and potential trend reversals.
//
//	PCDS = Ema(13, Ema(25, (Current - Prior)))
//	APCDS = Ema(13, Ema(25, Abs(Current - Prior)))
//	TSI = (PCDS / APCDS) * 100
//
// APCDS (the smoothed absolute price change) is zero only when price has
// been perfectly flat, in which case there is no momentum to report; TSI is
// defined as 0, its neutral center on the signed -100 to 100 scale, instead
// of propagating a 0/0 NaN.
//
// Example:
//
//	tsi := trend.NewTsi[float64]()
//	result := tsi.Compute(closings)
type Tsi[T helper.Number] struct {
	// FirstSmoothing is the first smoothing moving average.
	FirstSmoothing Ma[T]

	// SecondSmoothing is the second smoothing moving average.
	SecondSmoothing Ma[T]
}

// NewTsi function initializes a new TSI instance with the default parameters.
func NewTsi[T helper.Float]() *Tsi[T] {
	return NewTsiWith[T](
		DefaultTsiFirstSmoothingPeriod,
		DefaultTsiSecondSmoothingPeriod,
	)
}

// NewTsiWith function initializes a new TSI instance with the given parameters.
func NewTsiWith[T helper.Float](firstSmoothingPeriod, secondSmoothingPeriod int) *Tsi[T] {
	return &Tsi[T]{
		FirstSmoothing:  NewEmaWithPeriod[T](firstSmoothingPeriod),
		SecondSmoothing: NewEmaWithPeriod[T](secondSmoothingPeriod),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the TSI over the specified period, supporting context cancellation.
func (t *Tsi[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	// Price change
	pcsSplice := helper.DuplicateWithContext(ctx, helper.ChangeWithContext(ctx, closings, 1),
		2,
	)

	//	PCDS = Ema(13, Ema(25, (Current - Prior)))
	pcds := ComputeMaWithContext(ctx, t.SecondSmoothing, ComputeMaWithContext(ctx, t.FirstSmoothing, pcsSplice[0]))

	// APCDS = Ema(13, Ema(25, Abs(Current - Prior)))
	apcds := ComputeMaWithContext(ctx, t.SecondSmoothing, ComputeMaWithContext(ctx, t.FirstSmoothing, helper.AbsWithContext(ctx, pcsSplice[1])))

	// TSI = (PCDS / APCDS) * 100
	tsi := helper.OperateWithContext(ctx, pcds, apcds, func(pcd, apcd T) T {
		return helper.SafeDivide(pcd, apcd, T(0)) * T(100)
	})

	return tsi
}

// IdlePeriod is the initial period that TSI yield any results.
func (t *Tsi[T]) IdlePeriod() int {
	return t.FirstSmoothing.IdlePeriod() + t.SecondSmoothing.IdlePeriod() + 1
}

// String is the string representation of the TSI.
func (t *Tsi[T]) String() string {
	return fmt.Sprintf("TSI(%s,%s)",
		t.FirstSmoothing.String(),
		t.SecondSmoothing.String(),
	)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (t *Tsi[T]) Compute(closings <-chan T) <-chan T {
	return t.ComputeWithContext(context.Background(), closings)
}
