// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultMfiPeriod is the default period of the MFI.
	DefaultMfiPeriod = 14
)

// Mfi holds configuration parameters for calculating the Money Flow Index (MFI). It analyzes both the closing price
// and the volume to measure to identify overbought and oversold states. It is similar to the Relative Strength
// Index (RSI), but it also uses the volume.
//
//	Raw Money Flow = Typical Price * Volume
//	Money Ratio = Positive Money Flow / Negative Money Flow
//	Money Flow Index = 100 - (100 / (1 + Money Ratio))
//
// Example:
//
//	mfi := volume.NewMfi[float64]()
//	result := mfi.Compute(highs, lows, closings, volumes)
type Mfi[T helper.Number] struct {
	// TypicalPrice is the Typical Price instance.
	TypicalPrice *trend.TypicalPrice[T]

	// Sum is the Moving Sum instance.
	Sum *trend.MovingSum[T]
}

// NewMfi function initializes a new MFI instance with the default parameters.
func NewMfi[T helper.Number]() *Mfi[T] {
	return NewMfiWithPeriod[T](DefaultMfiPeriod)
}

// NewMfiWithPeriod function initializes a new MFI instance with the given period.
func NewMfiWithPeriod[T helper.Number](period int) *Mfi[T] {
	return &Mfi[T]{
		TypicalPrice: trend.NewTypicalPrice[T](),
		Sum:          trend.NewMovingSumWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the MFI.
func (m *Mfi[T]) ComputeWithContext(ctx context.Context, highs, lows, closings, volumes <-chan T) <-chan T {
	//	Raw Money Flow = Typical Price * Volume
	rawMoneyFlowSplice := helper.DuplicateWithContext(ctx, helper.MultiplyWithContext(ctx, m.TypicalPrice.ComputeWithContext(ctx, highs, lows, closings),
		volumes,
	),
		2,
	)

	moneyFlowSplice := helper.DuplicateWithContext(ctx, helper.MultiplyWithContext(ctx, helper.SignWithContext(ctx, helper.ChangeWithContext(ctx, rawMoneyFlowSplice[0], 1)),
		helper.SkipWithContext(ctx, rawMoneyFlowSplice[1], 1),
	),
		2,
	)

	// Money Ratio = Positive Money Flow / Negative Money Flow
	moneyRatio := helper.DivideWithContext(ctx, m.Sum.ComputeWithContext(ctx, helper.KeepPositivesWithContext(ctx, moneyFlowSplice[0])),
		m.Sum.ComputeWithContext(ctx, helper.MultiplyByWithContext(ctx, helper.KeepNegativesWithContext(ctx, moneyFlowSplice[1]),
			-1,
		),
		),
	)

	// Money Flow Index = 100 - (100 / (1 + Money Ratio))
	return helper.IncrementByWithContext(ctx, helper.MultiplyByWithContext(ctx, helper.PowWithContext(ctx, helper.IncrementByWithContext(ctx, moneyRatio,
		1,
	),
		-1,
	),
		-100,
	),
		100,
	)
}

// IdlePeriod is the initial period that MFI won't yield any results.
func (m *Mfi[T]) IdlePeriod() int {
	return m.Sum.IdlePeriod() + 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *Mfi[T]) Compute(highs, lows, closings, volumes <-chan T) <-chan T {
	return m.ComputeWithContext(context.Background(), highs, lows, closings, volumes)
}
