// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultMassIndexPeriod1 is the period for the first EMA.
	DefaultMassIndexPeriod1 = 9

	// DefaultMassIndexPeriod2 is the period for the second EMA.
	DefaultMassIndexPeriod2 = 9

	// DefaultMassIndexPeriod3 is the period for the third MovingSum.
	DefaultMassIndexPeriod3 = 25
)

// MassIndex represents the configuration parameters for calculating the
// Mass Index. It uses the high-low range to identify trend reversals
// based on range expansions.
//
//	Single EMA = EMA(9, Highs - Lows)
//	Double EMA = EMA(9, Single EMA)
//	Ratio = Single EMA / Double Ema1
//	Mass Index = SUM(Ratio, 25)
//
// Example:
type MassIndex[T helper.Number] struct {
	Ema1      *Ema[T]
	Ema2      *Ema[T]
	MovingSum *MovingSum[T]
}

// NewMassIndex function initializes a new APO instance
// with the default parameters.
func NewMassIndex[T helper.Number]() *MassIndex[T] {
	mi := &MassIndex[T]{
		Ema1:      NewEma[T](),
		Ema2:      NewEma[T](),
		MovingSum: NewMovingSum[T](),
	}

	mi.Ema1.Period = DefaultMassIndexPeriod1
	mi.Ema2.Period = DefaultMassIndexPeriod2
	mi.MovingSum.Period = DefaultMassIndexPeriod3

	return mi
}

// ComputeWithContext function takes a channel of numbers and computes the Mass Index.
func (m *MassIndex[T]) ComputeWithContext(ctx context.Context, highs, lows <-chan T) <-chan T {
	ema1 := helper.DuplicateWithContext(ctx, m.Ema1.ComputeWithContext(ctx, helper.SubtractWithContext(ctx, highs, lows)),
		2,
	)

	ema2 := m.Ema2.ComputeWithContext(ctx, ema1[0])
	ema1[1] = helper.SkipWithContext(ctx, ema1[1], m.Ema2.Period-1)

	ratio := helper.DivideWithContext(ctx, ema1[1], ema2)
	mi := m.MovingSum.ComputeWithContext(ctx, ratio)

	return mi
}

// IdlePeriod is the initial period that Mass Index won't yield any results.
func (m *MassIndex[T]) IdlePeriod() int {
	return m.Ema1.Period + m.Ema2.Period + m.MovingSum.Period - 3
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *MassIndex[T]) Compute(highs, lows <-chan T) <-chan T {
	return m.ComputeWithContext(context.Background(), highs, lows)
}
