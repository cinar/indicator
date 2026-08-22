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
//	%K1, %D1 = Stochastic(MACD, kPeriod, dPeriod)
//	%K2, %D2 = Stochastic(%D1, kPeriod, dPeriod)
//
//	STC = %D2
//
// The Stochastic pass (rolling-min/max normalization to a 0-100 range,
// then smoothed by an SMA) is applied twice, once to MACD and again to
// the first pass's %D, the way the standard Schaff Trend Cycle algorithm
// double-smooths MACD -- not once to MACD with a second division against
// its own %K/%D, which isn't bounded to 0-100 and can divide by a
// near-zero denominator.
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

	// %K1 and %K2 are unused (STC only needs the doubly-smoothed %D), but
	// each comes from the same internal duplicate fan-out as %D1/%D2, so
	// it still has to be drained -- an unread duplicate branch blocks the
	// shared producer, stalling %D1/%D2 too.
	k1, d1 := s.Stochastic.ComputeWithContext(ctx, macd)
	go helper.DrainWithContext(ctx, k1)

	k2, d2 := s.Stochastic.ComputeWithContext(ctx, d1)
	go helper.DrainWithContext(ctx, k2)

	return d2
}

// IdlePeriod is the initial period that STC won't yield any results.
func (s *Stc[T]) IdlePeriod() int {
	return s.Apo.IdlePeriod() + 2*s.Stochastic.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *Stc[T]) Compute(c <-chan T) <-chan T { return s.ComputeWithContext(context.Background(), c) }
