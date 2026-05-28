// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultRsiPeriod is the default period for the Relative Strength Index (RSI).
	DefaultRsiPeriod = 14
)

// Rsi represents the configuration parameter for calculating the Relative Strength Index (RSI).  It is a momentum
// indicator that measures the magnitude of recent price changes to evaluate overbought and oversold conditions.
//
//	RS = Average Gain / Average Loss
//	RSI = 100 - (100 / (1 + RS))
//
// Example:
//
//	rsi := momentum.NewRsi[float64]()
//	result := rsi.Compute(closings)
type Rsi[T helper.Number] struct {
	// Rma is the RMA instance.
	Rma *trend.Rma[T]
}

// NewRsi function initializes a new Relative Strength Index instance with the default parameters.
func NewRsi[T helper.Number]() *Rsi[T] {
	return NewRsiWithPeriod[T](DefaultRsiPeriod)
}

// NewRsiWithPeriod function initializes a new Relative Strength Index instance with the given period.
func NewRsiWithPeriod[T helper.Number](period int) *Rsi[T] {
	return &Rsi[T]{
		Rma: trend.NewRmaWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of closings numbers and computes the Relative Strength Index.
func (r *Rsi[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	changesSplice := helper.DuplicateWithContext(ctx, helper.ChangeWithContext(ctx, closings, 1),
		2,
	)

	averageGains := r.Rma.ComputeWithContext(ctx, helper.KeepPositivesWithContext(ctx, changesSplice[0]))

	averageLosses := helper.MultiplyByWithContext(ctx, r.Rma.ComputeWithContext(ctx, helper.KeepNegativesWithContext(ctx, changesSplice[1])),
		-1,
	)

	rs := helper.DivideWithContext(ctx, averageGains,
		averageLosses,
	)

	// RSI = 100 - (100 / (1 + RS))
	rsi := helper.IncrementByWithContext(ctx, // - (100 / (1 + RS))
		helper.MultiplyByWithContext(ctx, // 100 / (1 + RS)
			helper.MultiplyByWithContext(ctx, // 1 / (1 + RS)
				helper.PowWithContext(ctx, // 1 + RS
					helper.IncrementByWithContext(ctx, rs, 1),
					-1,
				),
				100,
			),
			-1,
		),
		100,
	)

	return rsi
}

// IdlePeriod is the initial period that Relative Strength Index won't yield any results.
func (r *Rsi[T]) IdlePeriod() int {
	return r.Rma.IdlePeriod() + 1
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (r *Rsi[T]) Compute(closings <-chan T) <-chan T {
	return r.ComputeWithContext(context.Background(), closings)
}
