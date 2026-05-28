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
	// DefaultStochasticRsiPeriod is the default period for the Stochastic Relative Strength Index (RSI).
	DefaultStochasticRsiPeriod = 14
)

// StochasticRsi represents the configuration parameter for calculating the Stochastic Relative Strength Index (RSI).
// It is a momentum indicator that focuses on the historical performance to evaluate overbought and
// oversold conditions.
//
//	                    RSI - Min(RSI)
//	Stochastic RSI = -------------------------
//	                   Max(RSI) - Min(RSI)
//
// Example:
//
//	stochasticRsi := momentum.NewStochasticRsi[float64]()
//	result := stochasticRsi.Compute(closings)
type StochasticRsi[T helper.Number] struct {
	// Rsi is that RSI instance.
	Rsi *Rsi[T]

	// Min is the Moving Min instance.
	Min *trend.MovingMin[T]

	// Max is the Moving Max instance.
	Max *trend.MovingMax[T]
}

// NewStochasticRsi function initializes a new Storchastic RSI instance with the default parameters.
func NewStochasticRsi[T helper.Number]() *StochasticRsi[T] {
	return NewStochasticRsiWithPeriod[T](DefaultStochasticRsiPeriod)
}

// NewStochasticRsiWithPeriod function initializes a new Stochastic RSI instance with the given period.
func NewStochasticRsiWithPeriod[T helper.Number](period int) *StochasticRsi[T] {
	return &StochasticRsi[T]{
		Rsi: NewRsiWithPeriod[T](period),
		Min: trend.NewMovingMinWithPeriod[T](period),
		Max: trend.NewMovingMaxWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of closings numbers and computes the Stochastic RSI.
func (s *StochasticRsi[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	rsisSplice := helper.DuplicateWithContext(ctx, s.Rsi.ComputeWithContext(ctx, closings),
		3,
	)

	rsisSplice[0] = helper.SkipWithContext(ctx, rsisSplice[0], s.Max.IdlePeriod())

	minRsisSplice := helper.DuplicateWithContext(ctx, s.Min.ComputeWithContext(ctx, rsisSplice[1]),
		2,
	)

	maxRsis := s.Max.ComputeWithContext(ctx, rsisSplice[2])

	result := helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, rsisSplice[0],
		minRsisSplice[0],
	),
		helper.SubtractWithContext(ctx, maxRsis,
			minRsisSplice[1],
		),
	)

	return result
}

// IdlePeriod is the initial period that Stochasic RSI won't yield any results.
func (s *StochasticRsi[T]) IdlePeriod() int {
	return s.Rsi.IdlePeriod() + s.Min.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *StochasticRsi[T]) Compute(closings <-chan T) <-chan T {
	return s.ComputeWithContext(context.Background(), closings)
}
