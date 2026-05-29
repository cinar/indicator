// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"math"

	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultUltimateOscillatorShortPeriod is the default short period for the Ultimate Oscillator (UO).
	DefaultUltimateOscillatorShortPeriod = 7

	// DefaultUltimateOscillatorMediumPeriod is the default medium period for the Ultimate Oscillator (UO).
	DefaultUltimateOscillatorMediumPeriod = 14

	// DefaultUltimateOscillatorLongPeriod is the default long period for the Ultimate Oscillator (UO).
	DefaultUltimateOscillatorLongPeriod = 28
)

// UltimateOscillator represents the configuration parameter for calculating the Ultimate Oscillator (UO).
// It was developed by Larry Williams in 1976 to measure the price momentum of an asset across multiple
// timeframes. By using the weighted average of three different timeframes the indicator has less
// volatility and fewer trade signals compared to other oscillators that rely on a single timeframe.
//
//	BP = Close - Minimum(Low, Prior Close)
//	TR = Maximum(High, Prior Close) - Minimum(Low, Prior Close)
//	Average7 = Sum(BP for 7 periods) / Sum(TR for 7 periods)
//	Average14 = Sum(BP for 14 periods) / Sum(TR for 14 periods)
//	Average28 = Sum(BP for 28 periods) / Sum(TR for 28 periods)
//	UO = 100 * [(4 * Average7) + (2 * Average14) + Average28] / (4 + 2 + 1)
//
// Example:
//
//	uo := momentum.NewUltimateOscillator[float64]()
//	values := uo.Compute(highs, lows, closings)
type UltimateOscillator[T helper.Number] struct {
	// ShortPeriod is the short period for the UO.
	ShortPeriod int

	// MediumPeriod is the medium period for the UO.
	MediumPeriod int

	// LongPeriod is the long period for the UO.
	LongPeriod int
}

// NewUltimateOscillator function initializes a new Ultimate Oscillator instance.
func NewUltimateOscillator[T helper.Number]() *UltimateOscillator[T] {
	return NewUltimateOscillatorWithPeriods[T](
		DefaultUltimateOscillatorShortPeriod,
		DefaultUltimateOscillatorMediumPeriod,
		DefaultUltimateOscillatorLongPeriod,
	)
}

// NewUltimateOscillatorWithPeriods function initializes a new Ultimate Oscillator instance with the given periods.
func NewUltimateOscillatorWithPeriods[T helper.Number](shortPeriod, mediumPeriod, longPeriod int) *UltimateOscillator[T] {
	return &UltimateOscillator[T]{
		ShortPeriod:  shortPeriod,
		MediumPeriod: mediumPeriod,
		LongPeriod:   longPeriod,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Ultimate Oscillator.
func (u *UltimateOscillator[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)
	highsSplice := helper.DuplicateWithContext(ctx, highs, 1)
	lowsSplice := helper.DuplicateWithContext(ctx, lows, 1)

	priorClosings := helper.BufferedWithContext(ctx, closingsSplice[0], 1)
	currentClosings := helper.SkipWithContext(ctx, helper.BufferedWithContext(ctx, closingsSplice[1], 1), 1)
	highs = helper.SkipWithContext(ctx, helper.BufferedWithContext(ctx, highsSplice[0], 1), 1)
	lows = helper.SkipWithContext(ctx, helper.BufferedWithContext(ctx, lowsSplice[0], 1), 1)

	type bpTr struct {
		bp T
		tr T
	}

	bpTrChan := helper.Operate4WithContext(ctx, currentClosings, highs, lows, priorClosings, func(c, h, l, pc T) bpTr {
		minLowPc := T(math.Min(float64(l), float64(pc)))
		maxHighPc := T(math.Max(float64(h), float64(pc)))

		return bpTr{
			bp: c - minLowPc,
			tr: maxHighPc - minLowPc,
		}
	})

	bpTrSplice := helper.DuplicateWithContext(ctx, bpTrChan, 2)

	bpChan := helper.MapWithContext(ctx, helper.BufferedWithContext(ctx, bpTrSplice[0], 1), func(bt bpTr) T {
		return bt.bp
	})

	trChan := helper.MapWithContext(ctx, helper.BufferedWithContext(ctx, bpTrSplice[1], 1), func(bt bpTr) T {
		return bt.tr
	})

	bpSpliced := helper.DuplicateWithContext(ctx, bpChan, 3)
	trSpliced := helper.DuplicateWithContext(ctx, trChan, 3)

	shortBpSum := trend.NewMovingSumWithPeriod[T](u.ShortPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, bpSpliced[0], 1))
	shortTrSum := trend.NewMovingSumWithPeriod[T](u.ShortPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, trSpliced[0], 1))

	mediumBpSum := trend.NewMovingSumWithPeriod[T](u.MediumPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, bpSpliced[1], 1))
	mediumTrSum := trend.NewMovingSumWithPeriod[T](u.MediumPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, trSpliced[1], 1))

	longBpSum := trend.NewMovingSumWithPeriod[T](u.LongPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, bpSpliced[2], 1))
	longTrSum := trend.NewMovingSumWithPeriod[T](u.LongPeriod).ComputeWithContext(ctx, helper.BufferedWithContext(ctx, trSpliced[2], 1))

	// Align sums to the long period
	shortBpSum = helper.SkipWithContext(ctx, shortBpSum, u.LongPeriod-u.ShortPeriod)
	shortTrSum = helper.SkipWithContext(ctx, shortTrSum, u.LongPeriod-u.ShortPeriod)
	mediumBpSum = helper.SkipWithContext(ctx, mediumBpSum, u.LongPeriod-u.MediumPeriod)
	mediumTrSum = helper.SkipWithContext(ctx, mediumTrSum, u.LongPeriod-u.MediumPeriod)

	avgShort := helper.DivideWithContext(ctx, shortBpSum, shortTrSum)
	avgMedium := helper.DivideWithContext(ctx, mediumBpSum, mediumTrSum)
	avgLong := helper.DivideWithContext(ctx, longBpSum, longTrSum)

	// UO = 100 * [(4 * Average7) + (2 * Average14) + Average28] / (4 + 2 + 1)
	// (4 + 2 + 1) = 7

	termShort := helper.MultiplyByWithContext(ctx, avgShort, 4)
	termMedium := helper.MultiplyByWithContext(ctx, avgMedium, 2)

	sumTerms := helper.AddWithContext(ctx, helper.AddWithContext(ctx, termShort, termMedium), avgLong)

	return helper.DivideByWithContext(ctx, helper.MultiplyByWithContext(ctx, sumTerms, 100), 7)
}

// IdlePeriod is the initial period that Ultimate Oscillator won't yield any results.
func (u *UltimateOscillator[T]) IdlePeriod() int {
	return u.LongPeriod
}

func (u *UltimateOscillator[T]) String() string {
	return "Ultimate Oscillator"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (u *UltimateOscillator[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return u.ComputeWithContext(context.Background(), highs, lows, closings)
}
