// Pvopyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultPvoShortPeriod is the default short period for the Percentage Volume Oscillator.
	DefaultPvoShortPeriod = 12

	// DefaultPvoLongPeriod is the default long period for the Percentage Volume Oscillator.
	DefaultPvoLongPeriod = 26

	// DefaultPvoSignalPeriod is the default signal period for the Percentage Volume Oscillator.
	DefaultPvoSignalPeriod = 9
)

// Pvo represents the configuration parameter for calculating the Percentage Volume Oscillator (PVO). It is a
// momentum oscillator for the price. It is used to indicate the ups and downs based on the price. A breakout
// is confirmed when PVO is positive.
//
//	PVO = ((EMA(shortPeriod, prices) - EMA(longPeriod, prices)) / EMA(longPeriod, prices)) * 100
//	Signal = EMA(9, PVO)
//	Histogram = PVO - Signal
//
// Example:
//
//	pvo := momentum.Pvo[float64]()
//	p, s, h := pvo.Compute(volumes)
type Pvo[T helper.Number] struct {
	// ShortEma is the short EMA instance.
	ShortEma *trend.Ema[T]

	// LongEma is the long EMA instance.
	LongEma *trend.Ema[T]

	// SignalEma is the signal EMA instance.
	SignalEma *trend.Ema[T]
}

// NewPvo function initializes a new Percentage Volume Oscillator instance.
func NewPvo[T helper.Number]() *Pvo[T] {
	return &Pvo[T]{
		ShortEma:  trend.NewEmaWithPeriod[T](DefaultPvoShortPeriod),
		LongEma:   trend.NewEmaWithPeriod[T](DefaultPvoLongPeriod),
		SignalEma: trend.NewEmaWithPeriod[T](DefaultPvoSignalPeriod),
	}
}

// Compute function takes a channel of numbers and computes the Percentage Volume Oscillator.
// Returns pvo, signal, histogram.
func (p *Pvo[T]) Compute(volumes <-chan T) (<-chan T, <-chan T, <-chan T) {
	volumesSplice := helper.Duplicate(
		volumes,
		2,
	)

	shortEma := p.ShortEma.Compute(volumesSplice[0])

	longEmaSplice := helper.Duplicate(
		p.LongEma.Compute(volumesSplice[1]),
		2,
	)

	shortEma = helper.Skip(shortEma, p.LongEma.IdlePeriod()-p.ShortEma.IdlePeriod())

	//	PVO = ((EMA(shortPeriod, prices) - EMA(longPeriod, prices)) / EMA(longPeriod, prices)) * 100
	pvoSplice := helper.Duplicate(
		helper.MultiplyBy(
			helper.Divide(
				helper.Subtract(shortEma, longEmaSplice[0]),
				longEmaSplice[1],
			),
			100,
		),
		3,
	)

	//	Signal = EMA(9, PVO)
	signalSplice := helper.Duplicate(
		p.SignalEma.Compute(pvoSplice[0]),
		2,
	)

	pvoSplice[1] = helper.Skip(pvoSplice[1], p.SignalEma.IdlePeriod())
	pvoSplice[2] = helper.Skip(pvoSplice[2], p.SignalEma.IdlePeriod())

	//	Histogram = PVO - Signal
	histogram := helper.Subtract(pvoSplice[1], signalSplice[0])

	return pvoSplice[2], signalSplice[1], histogram
}

// IdlePeriod is the initial period that Percentage Volume Oscillator won't yield any results.
func (p *Pvo[T]) IdlePeriod() int {
	return p.LongEma.IdlePeriod() + p.SignalEma.IdlePeriod()
}
