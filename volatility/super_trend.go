// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultSuperTrendPeriod is the default period value.
	DefaultSuperTrendPeriod = 14

	// DefaultSuperTrendMultiplier is the default multiplier value.
	DefaultSuperTrendMultiplier = 2.5
)

// SuperTrend represents the configuration parameters for calculating the Super Trend.
//
//	BasicUpperBands = (High + Low) / 2 + Multiplier * ATR
//	BasicLowerBands = (High + Low) / 2 - Multiplier * ATR
//	FinalUpperBands = If (BasicUpperBand < PreviousFinalUpperBand)
//	                  Or (PreviousClose > PreviousFinalUpperBand)
//	                  Then BasicUpperBand Else PreviousFinalUpperBand
//	FinalLowerBands = If (BasicLowerBand > PreviousFinalLowerBand)
//	                  Or (PreviousClose < PreviousFinalLowerBand)
//	                  Then BasicLowerBand Else PreviousFinalLowerBand
//	SuperTrend = If upTrend
//				 Then
//	               If (Close <= FinalUpperBand) Then FinalUpperBand Else FinalLowerBand
//	             Else
//	               If (Close >= FinalLowerBand) Then FinalLowerBand Else FinalUpperBand
//
//	UpTrend = If (SuperTrend == FinalUpperBand) Then True Else False
//
// Example:
type SuperTrend[T helper.Number] struct {
	Atr        *Atr[T]
	Multiplier T
}

// NewSuperTrend function initializes a new Super Trend instance with the default parameters.
func NewSuperTrend[T helper.Number]() *SuperTrend[T] {
	multiplier := DefaultSuperTrendMultiplier

	return NewSuperTrendWithPeriod[T](
		DefaultSuperTrendPeriod,
		T(multiplier),
	)
}

// NewSuperTrendWithPeriod initializes a new Super Trend instance with the given period and multiplier.
func NewSuperTrendWithPeriod[T helper.Number](period int, multiplier T) *SuperTrend[T] {
	return NewSuperTrendWithMa[T](
		trend.NewHmaWithPeriod[T](period),
		multiplier,
	)
}

// NewSuperTrendWithMa function initializes a new Super Trend instance with the given moving average instance
// and multiplier.
func NewSuperTrendWithMa[T helper.Number](ma trend.Ma[T], multiplier T) *SuperTrend[T] {
	return &SuperTrend[T]{
		Atr:        NewAtrWithMa(ma),
		Multiplier: multiplier,
	}
}

// Compute function calculates the Super Trend, using separate channels for highs, lows, and closings.
func (s *SuperTrend[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)
	closingsSplice := helper.Duplicate(closings, 2)

	medians :=
		helper.Skip(
			helper.DivideBy(
				helper.Add(highsSplice[0], lowsSplice[0]),
				2,
			),
			s.Atr.IdlePeriod(),
		)

	atrMultiples :=
		helper.MultiplyBy(
			s.Atr.Compute(highsSplice[1], lowsSplice[1], closingsSplice[0]),
			s.Multiplier,
		)

	closingsSplice[1] = helper.Skip(closingsSplice[1], s.Atr.IdlePeriod())

	first := true
	upTrend := false
	var previousClosing T
	var finalUpperBand T
	var finalLowerBand T

	superTrend := helper.Operate3(medians, atrMultiples, closingsSplice[1], func(median, atrMultiple, closing T) T {
		//	BasicUpperBands = (High + Low) / 2 + Multiplier * ATR
		basicUpperBand := median + atrMultiple

		//	BasicLowerBands = (High + Low) / 2 - Multiplier * ATR
		basicLowerBand := median - atrMultiple

		var superTrend T

		if first {
			first = false
			finalUpperBand = basicUpperBand
			finalLowerBand = basicLowerBand
			superTrend = finalLowerBand
		} else {
			//	FinalUpperBands = If (BasicUpperBand < PreviousFinalUpperBand)
			//	                  Or (PreviousClose > PreviousFinalUpperBand)
			//	                  Then BasicUpperBand Else PreviousFinalUpperBand
			if (basicUpperBand < finalUpperBand) || (previousClosing > finalUpperBand) {
				finalUpperBand = basicUpperBand
			}

			//	FinalLowerBands = If (BasicLowerBand > PreviousFinalLowerBand)
			//	                  Or (PreviousClose < PreviousFinalLowerBand)
			//	                  Then BasicLowerBand Else PreviousFinalLowerBand
			if (basicLowerBand > finalLowerBand) || (previousClosing < finalLowerBand) {
				finalLowerBand = basicLowerBand
			}

			//	SuperTrend = If upTrend
			//				 Then
			//	               If (Close <= FinalUpperBand) Then FinalUpperBand Else FinalLowerBand
			//	             Else
			//	               If (Close >= FinalLowerBand) Then FinalLowerBand Else FinalUpperBand
			//
			//	UpTrend = If (SuperTrend == FinalUpperBand) Then True Else False
			if upTrend {
				if closing <= finalUpperBand {
					superTrend = finalUpperBand
				} else {
					superTrend = finalLowerBand
					upTrend = false
				}
			} else {
				if closing >= finalLowerBand {
					superTrend = finalLowerBand
				} else {
					superTrend = finalUpperBand
					upTrend = true
				}
			}
		}

		previousClosing = closing

		return superTrend
	})

	return superTrend
}

// IdlePeriod is the initial period that Super Trend won't yield any results.
func (s *SuperTrend[T]) IdlePeriod() int {
	return s.Atr.IdlePeriod()
}
