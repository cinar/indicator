// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultStcFastPeriod is the default fast EMA period for STC.
	DefaultStcFastPeriod = 23

	// DefaultStcSlowPeriod is the default slow EMA period for STC.
	DefaultStcSlowPeriod = 50

	// DefaultStcKPeriod is the default window length of each stochastic pass.
	DefaultStcKPeriod = 10

	// DefaultStcDPeriod is the default smoothing period; 3 gives a factor of 0.5.
	DefaultStcDPeriod = 3
)

// Stc represents the configuration parameters for calculating the
// Schaff Trend Cycle (STC) indicator. It applies a stochastic
// normalization to MACD twice, smoothing the result of each pass
// exponentially, as in Doug Schaff's definition and TA-Lib STC.
//
//	MACD = EMA(values, fastPeriod) - EMA(values, slowPeriod)
//
//	%K1 = 100 * (MACD - Min(MACD, kPeriod)) / (Max(MACD, kPeriod) - Min(MACD, kPeriod))
//	PF  = PF[previous] + factor * (%K1 - PF[previous])
//
//	%K2 = 100 * (PF - Min(PF, kPeriod)) / (Max(PF, kPeriod) - Min(PF, kPeriod))
//	STC = STC[previous] + factor * (%K2 - STC[previous])
//
//	factor = 2 / (dPeriod + 1)
//
// Both EMAs start on the same bar, the slow EMA's first: the fast EMA is
// seeded with the average of the fastPeriod values before it, rather than
// of the first fastPeriod values, as in TA-Lib STC and MACD. This differs
// from [Apo] and [Macd], whose fast EMA starts earlier; the two converge
// as the series grows.
//
// The default dPeriod of 3 gives Schaff's factor of 0.5. Each smoothed
// series starts from its first %K value. When the window's maximum equals
// its minimum, %K is undefined and the previous smoothed value is kept
// (0 if there is none yet).
//
// Example:
//
//	stc := trend.NewStc[float64]()
//	result := stc.Compute(closings)
type Stc[T helper.Float] struct {
	// FastPeriod is the period for the fast EMA.
	FastPeriod int

	// SlowPeriod is the period for the slow EMA.
	SlowPeriod int

	// KPeriod is the window length of each stochastic pass (the cycle length).
	KPeriod int

	// DPeriod sets the smoothing factor, 2 / (DPeriod + 1), applied after
	// each stochastic pass.
	DPeriod int

	// FastEma is the fast EMA of the MACD.
	FastEma *Ema[T]

	// SlowEma is the slow EMA of the MACD.
	SlowEma *Ema[T]
}

// NewStc function initializes a new STC instance with the default parameters.
func NewStc[T helper.Float]() *Stc[T] {
	return NewStcWithPeriod[T](
		DefaultStcFastPeriod,
		DefaultStcSlowPeriod,
		DefaultStcKPeriod,
		DefaultStcDPeriod,
	)
}

// NewStcWithPeriod function initializes a new STC instance with the given periods.
func NewStcWithPeriod[T helper.Float](fastPeriod, slowPeriod, kPeriod, dPeriod int) *Stc[T] {
	return &Stc[T]{
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
		KPeriod:    kPeriod,
		DPeriod:    dPeriod,
		FastEma:    NewEmaWithPeriod[T](fastPeriod),
		SlowEma:    NewEmaWithPeriod[T](slowPeriod),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the STC indicator.
func (s *Stc[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	c = helper.BufferedWithContext(ctx, c, s.SlowPeriod)
	cs := helper.DuplicateWithContext(ctx, c, 2)

	// Skipping the first SlowPeriod-FastPeriod values makes the fast EMA
	// start on the same bar as the slow EMA.
	fast := s.FastEma.ComputeWithContext(ctx, helper.SkipWithContext(ctx, cs[0], s.SlowPeriod-s.FastPeriod))
	slow := s.SlowEma.ComputeWithContext(ctx, cs[1])
	macd := helper.SubtractWithContext(ctx, fast, slow)

	pf := s.smoothedStochasticWithContext(ctx, macd)

	return s.smoothedStochasticWithContext(ctx, pf)
}

// smoothedStochasticWithContext normalizes the values to 0-100 against the
// minimum and maximum of the last KPeriod values, and smooths the result
// exponentially.
func (s *Stc[T]) smoothedStochasticWithContext(ctx context.Context, values <-chan T) <-chan T {
	vs := helper.DuplicateWithContext(ctx, values, 3)

	maxs := NewMovingMaxWithPeriod[T](s.KPeriod)
	mins := NewMovingMinWithPeriod[T](s.KPeriod)

	highest := maxs.ComputeWithContext(ctx, vs[0])
	lowest := mins.ComputeWithContext(ctx, vs[1])
	current := helper.SkipWithContext(ctx, vs[2], maxs.IdlePeriod())

	factor := T(2) / T(s.DPeriod+1)
	started := false

	var smoothed T

	return helper.Operate3WithContext(ctx, current, highest, lowest, func(value, high, low T) T {
		k := smoothed
		if high > low {
			k = 100 * (value - low) / (high - low)
		}

		if started {
			smoothed += factor * (k - smoothed)
		} else {
			smoothed = k
			started = true
		}

		return smoothed
	})
}

// IdlePeriod is the initial period that STC won't yield any results.
func (s *Stc[T]) IdlePeriod() int {
	return s.SlowEma.IdlePeriod() + 2*(s.KPeriod-1)
}

// String is the string representation of the STC.
func (s *Stc[T]) String() string {
	return fmt.Sprintf("STC(%d,%d,%d,%d)", s.FastPeriod, s.SlowPeriod, s.KPeriod, s.DPeriod)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *Stc[T]) Compute(c <-chan T) <-chan T { return s.ComputeWithContext(context.Background(), c) }
