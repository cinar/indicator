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
	// DefaultPpoShortPeriod is the default short period for the Percentage Price Oscillator.
	DefaultPpoShortPeriod = 12

	// DefaultPpoLongPeriod is the default long period for the Percentage Price Oscillator.
	DefaultPpoLongPeriod = 26

	// DefaultPpoSignalPeriod is the default signal period for the Percentage Price Oscillator.
	DefaultPpoSignalPeriod = 9
)

// Ppo represents the configuration parameter for calculating the Percentage Price Oscillator (PPO). It is a momentum
// oscillator for the price. It is used to indicate the ups and downs based on the price. A breakout is confirmed
// when PPO is positive.
//
//	PPO = ((EMA(shortPeriod, prices) - EMA(longPeriod, prices)) / EMA(longPeriod, prices)) * 100
//	Signal = EMA(9, PPO)
//	Histogram = PPO - Signal
//
// Example:
//
//	ppo := momentum.Ppo[float64]()
//	p, s, h := ppo.Compute(closings)
type Ppo[T helper.Number] struct {
	// ShortEma is the short EMA instance.
	ShortEma *trend.Ema[T]

	// LongEma is the long EMA instance.
	LongEma *trend.Ema[T]

	// SignalEma is the signal EMA instance.
	SignalEma *trend.Ema[T]
}

// NewPpo function initializes a new Percentage Price Oscillator instance.
func NewPpo[T helper.Number]() *Ppo[T] {
	return &Ppo[T]{
		ShortEma:  trend.NewEmaWithPeriod[T](DefaultPpoShortPeriod),
		LongEma:   trend.NewEmaWithPeriod[T](DefaultPpoLongPeriod),
		SignalEma: trend.NewEmaWithPeriod[T](DefaultPpoSignalPeriod),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Percentage Price Oscillator.
// Returns ppo, signal, histogram.
func (p *Ppo[T]) ComputeWithContext(ctx context.Context, closings <-chan T) (<-chan T, <-chan T, <-chan T) {
	closingsSplice := helper.DuplicateWithContext(ctx, closings,
		2,
	)

	shortEma := p.ShortEma.ComputeWithContext(ctx, closingsSplice[0])

	longEmaSplice := helper.DuplicateWithContext(ctx, p.LongEma.ComputeWithContext(ctx, closingsSplice[1]),
		2,
	)

	shortEma = helper.SkipWithContext(ctx, shortEma, p.LongEma.IdlePeriod()-p.ShortEma.IdlePeriod())

	//	PPO = ((EMA(shortPeriod, prices) - EMA(longPeriod, prices)) / EMA(longPeriod, prices)) * 100
	ppoSplice := helper.DuplicateWithContext(ctx, helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, shortEma, longEmaSplice[0]),
		longEmaSplice[1],
	),
		100,
	),
		3,
	)

	//	Signal = EMA(9, PPO)
	signalSplice := helper.DuplicateWithContext(ctx, p.SignalEma.ComputeWithContext(ctx, ppoSplice[0]),
		2,
	)

	ppoSplice[1] = helper.SkipWithContext(ctx, ppoSplice[1], p.SignalEma.IdlePeriod())
	ppoSplice[2] = helper.SkipWithContext(ctx, ppoSplice[2], p.SignalEma.IdlePeriod())

	//	Histogram = PPO - Signal
	histogram := helper.SubtractWithContext(ctx, ppoSplice[1], signalSplice[0])

	return ppoSplice[2], signalSplice[1], histogram
}

// IdlePeriod is the initial period that Percentage Price Oscillator won't yield any results.
func (p *Ppo[T]) IdlePeriod() int {
	return p.LongEma.IdlePeriod() + p.SignalEma.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (p *Ppo[T]) Compute(closings <-chan T) (<-chan T, <-chan T, <-chan T) {
	return p.ComputeWithContext(context.Background(), closings)
}
