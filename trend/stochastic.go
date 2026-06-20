// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultStochasticPeriod is the default period for the Stochastic indicator.
	DefaultStochasticPeriod = 10

	// DefaultStochasticSmaPeriod is the default period for the SMA of %K.
	DefaultStochasticSmaPeriod = 3
)

// Stochastic represents the configuration parameters for calculating
// the Stochastic indicator on a single input series. This is different
// from the Stochastic Oscillator which operates on high, low, and close.
// This generic version is useful for applying stochastic calculation to
// any series, such as MACD values in the Schaff Trend Cycle (STC).
//
//	K = (Value - Min(Value, period)) / (Max(Value, period) - Min(Value, period)) * 100
//	D = SMA(K, dPeriod)
//
// Example:
//
//	s := trend.NewStochastic[float64]()
//	k, d := s.Compute(values)
type Stochastic[T helper.Number] struct {
	// Period is the period for the min/max calculation.
	Period int

	// Sma is the SMA instance for %D calculation.
	Sma *Sma[T]
}

// NewStochastic function initializes a new Stochastic instance with the default parameters.
func NewStochastic[T helper.Number]() *Stochastic[T] {
	return NewStochasticWithPeriod[T](DefaultStochasticPeriod)
}

// NewStochasticWithPeriod function initializes a new Stochastic instance with the given period.
func NewStochasticWithPeriod[T helper.Number](period int) *Stochastic[T] {
	return &Stochastic[T]{
		Period: period,
		Sma:    NewSmaWithPeriod[T](DefaultStochasticSmaPeriod),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Stochastic indicator.
// Returns %K and %D.
func (s *Stochastic[T]) ComputeWithContext(ctx context.Context, values <-chan T) (<-chan T, <-chan T) {
	movingMin := NewMovingMinWithPeriod[T](s.Period)
	movingMax := NewMovingMaxWithPeriod[T](s.Period)

	values = helper.BufferedWithContext(ctx, values, s.Period)
	inputs := helper.DuplicateWithContext(ctx, values, 3)

	lowestSplice := helper.DuplicateWithContext(ctx, movingMin.ComputeWithContext(ctx, inputs[0]),
		2,
	)

	highest := movingMax.ComputeWithContext(ctx, inputs[1])

	skipped := helper.SkipWithContext(ctx, inputs[2], movingMin.IdlePeriod())

	kSplice := helper.DuplicateWithContext(ctx, helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, skipped, lowestSplice[0]),
		helper.SubtractWithContext(ctx, highest, lowestSplice[1]),
	),
		100,
	),
		2,
	)

	d := s.Sma.ComputeWithContext(ctx, kSplice[0])
	kSplice[1] = helper.SkipWithContext(ctx, kSplice[1], s.Sma.IdlePeriod())

	return kSplice[1], d
}

// IdlePeriod is the initial period that Stochastic won't yield any results.
func (s *Stochastic[T]) IdlePeriod() int {
	return s.Period + s.Sma.Period - 2
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *Stochastic[T]) Compute(values <-chan T) (<-chan T, <-chan T) {
	return s.ComputeWithContext(context.Background(), values)
}
