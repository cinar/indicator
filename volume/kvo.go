// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultKvoShortPeriod is the default short period for the Klinger Volume Oscillator.
	DefaultKvoShortPeriod = 34

	// DefaultKvoLongPeriod is the default long period for the Klinger Volume Oscillator.
	DefaultKvoLongPeriod = 55

	// DefaultKvoSignalPeriod is the default signal period for the Klinger Volume Oscillator.
	DefaultKvoSignalPeriod = 13
)

// Kvo represents the configuration parameters for calculating the Klinger Volume Oscillator (KVO).
// It is a volume-based oscillator that identifies long-term money flow trends using EMA differences.
//
//	Trend = +1 if High > High[1] and Low >= Low[1]
//	Trend = -1 if High <= High[1] and Low < Low[1]
//	Trend = 0 otherwise
//	VF = Volume * Trend
//	KVO = EMA(VF, shortPeriod) - EMA(VF, longPeriod)
//	Signal = EMA(KVO, signalPeriod)
//
// Example:
//
//	kvo := volume.NewKvo[float64]()
//	kvoResult, signalResult := kvo.Compute(highs, lows, volumes)
type Kvo[T helper.Number] struct {
	// ShortEma is the short EMA instance.
	ShortEma *trend.Ema[T]

	// LongEma is the long EMA instance.
	LongEma *trend.Ema[T]

	// SignalEma is the signal EMA instance.
	SignalEma *trend.Ema[T]
}

// NewKvo function initializes a new KVO instance.
func NewKvo[T helper.Number]() *Kvo[T] {
	return &Kvo[T]{
		ShortEma:  trend.NewEmaWithPeriod[T](DefaultKvoShortPeriod),
		LongEma:   trend.NewEmaWithPeriod[T](DefaultKvoLongPeriod),
		SignalEma: trend.NewEmaWithPeriod[T](DefaultKvoSignalPeriod),
	}
}

// Compute function takes channels of numbers and computes the Klinger Volume Oscillator.
// Returns kvo and signal.
func (k *Kvo[T]) Compute(highs, lows, volumes <-chan T) (<-chan T, <-chan T) {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)

	previousHighs := helper.Shift(highsSplice[0], 1, 0)
	previousLows := helper.Shift(lowsSplice[0], 1, 0)

	highsCopy := highsSplice[1]
	lowsCopy := lowsSplice[1]

	vf := helper.Operate5(highsCopy, previousHighs, lowsCopy, previousLows, volumes, func(high, prevHigh, low, prevLow, volume T) T {
		var trend T

		if high > prevHigh && low >= prevLow {
			trend = 1
		} else if high <= prevHigh && low < prevLow {
			trend = -1
		} else {
			trend = 0
		}

		return volume * trend
	})

	vfSplice := helper.Duplicate(helper.Skip(vf, 1), 2)

	shortEma := k.ShortEma.Compute(vfSplice[0])
	longEma := k.LongEma.Compute(vfSplice[1])

	shortEma = helper.Skip(shortEma, k.LongEma.IdlePeriod()-k.ShortEma.IdlePeriod())

	kvo := helper.Subtract(shortEma, longEma)

	kvoSplice := helper.Duplicate(kvo, 2)

	signal := k.SignalEma.Compute(kvoSplice[0])
	kvoResult := helper.Skip(kvoSplice[1], k.SignalEma.IdlePeriod())

	return kvoResult, signal
}

// IdlePeriod is the initial period that KVO won't yield any results.
func (k *Kvo[T]) IdlePeriod() int {
	return k.LongEma.IdlePeriod() + k.SignalEma.IdlePeriod() + 1
}
