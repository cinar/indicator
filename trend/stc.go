// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultStcFastPeriod is the default fast EMA period for STC.
	DefaultStcFastPeriod = 23

	// DefaultStcSlowPeriod is the default slow EMA period for STC.
	DefaultStcSlowPeriod = 50

	// DefaultStcKPeriod is the default period for the Stochastic %K.
	DefaultStcKPeriod = 10

	// DefaultStcDPeriod is the default period for the Stochastic %D.
	DefaultStcDPeriod = 3
)

// Stc represents the configuration parameters for calculating the
// Schaff Trend Cycle (STC) indicator. It combines MACD with
// stochastic oscillators to identify trend direction and potential
// entry points.
//
//	EMA1 = EMA(values, fastPeriod)
//	EMA2 = EMA(values, slowPeriod)
//	MACD = EMA1 - EMA2
//
//	%K = Stochastic %K of MACD with kPeriod
//	%D = Stochastic %D of MACD with dPeriod
//
//	STC = 100 * (MACD - %K) / (%D - %K)
//
// Example:
//
//	stc := trend.NewStc[float64]()
//	result := stc.Compute(closings)
type Stc[T helper.Number] struct {
	// FastPeriod is the period for the fast EMA.
	FastPeriod int

	// SlowPeriod is the period for the slow EMA.
	SlowPeriod int

	// KPeriod is the period for the Stochastic %K.
	KPeriod int

	// DPeriod is the period for the Stochastic %D.
	DPeriod int

	// Apo is the APO instance for MACD calculation.
	Apo *Apo[T]

	// Stochastic is the Stochastic instance.
	Stochastic *Stochastic[T]
}

// NewStc function initializes a new STC instance with the default parameters.
func NewStc[T helper.Number]() *Stc[T] {
	return NewStcWithPeriod[T](
		DefaultStcFastPeriod,
		DefaultStcSlowPeriod,
		DefaultStcKPeriod,
		DefaultStcDPeriod,
	)
}

// NewStcWithPeriod function initializes a new STC instance with the given periods.
func NewStcWithPeriod[T helper.Number](fastPeriod, slowPeriod, kPeriod, dPeriod int) *Stc[T] {
	apo := NewApo[T]()
	apo.FastPeriod = fastPeriod
	apo.SlowPeriod = slowPeriod

	stochastic := NewStochasticWithPeriod[T](kPeriod)
	stochastic.Sma.Period = dPeriod

	return &Stc[T]{
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
		KPeriod:    kPeriod,
		DPeriod:    dPeriod,
		Apo:        apo,
		Stochastic: stochastic,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the STC indicator.
func (s *Stc[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	c = helper.BufferedWithContext(ctx, c, s.SlowPeriod)
	macd := s.Apo.ComputeWithContext(ctx, c)

	macd = helper.BufferedWithContext(ctx, macd, s.Stochastic.Period)
	inputs := helper.DuplicateWithContext(ctx, macd, 4)

	movingMin := NewMovingMinWithPeriod[T](s.Stochastic.Period)
	movingMax := NewMovingMaxWithPeriod[T](s.Stochastic.Period)

	lowestSplice := helper.DuplicateWithContext(ctx, movingMin.ComputeWithContext(ctx, inputs[0]),
		2,
	)

	highest := movingMax.ComputeWithContext(ctx, inputs[1])

	skipped := helper.SkipWithContext(ctx, inputs[2], movingMin.IdlePeriod())

	kValues := helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, skipped, lowestSplice[0]),
		helper.SubtractWithContext(ctx, highest, lowestSplice[1]),
	),
		100,
	)

	kDuplicate := helper.DuplicateWithContext(ctx, kValues, 2)

	d := s.Stochastic.Sma.ComputeWithContext(ctx, kDuplicate[0])

	kValues = helper.SkipWithContext(ctx, kDuplicate[1], s.Stochastic.Sma.IdlePeriod())

	macdForStc := helper.SkipWithContext(ctx, inputs[3], s.Stochastic.IdlePeriod())

	return helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, macdForStc, kValues),
		helper.SubtractWithContext(ctx, d, kValues),
	),
		100,
	)
}

// IdlePeriod is the initial period that STC won't yield any results.
func (s *Stc[T]) IdlePeriod() int {
	return s.Apo.IdlePeriod() + s.Stochastic.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *Stc[T]) Compute(c <-chan T) <-chan T { return s.ComputeWithContext(context.Background(), c) }
