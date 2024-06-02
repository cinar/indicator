// StochasticOscillatorpyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultStochasticOscillatorMaxAndMinPeriod is the default max and min period for the Stochastic Oscillator.
	DefaultStochasticOscillatorMaxAndMinPeriod = 14

	// DefaultStochasticOscillatorPeriod is the default period for the Stochastic Oscillator.
	DefaultStochasticOscillatorPeriod = 3
)

// StochasticOscillator represents the configuration parameter for calculating the Stochastic Oscillator. It is a
// momentum indicator that shows the location of the closing relative to high-low range over a set number
// of periods.
//
//	K = (Closing - Lowest Low) / (Highest High - Lowest Low) * 100
//	D = 3-Period SMA of K
//
// Example:
//
//	so := momentum.StochasticOscillator[float64]()
//	k, d := wr.Compute(highs, lows, closings)
type StochasticOscillator[T helper.Number] struct {
	// Max is the Moving Max instance.
	Max *trend.MovingMax[T]

	// Min is the Moving Min instance.
	Min *trend.MovingMin[T]

	// Sma is the SMA instance.
	Sma *trend.Sma[T]
}

// NewStochasticOscillator function initializes a new Stochastic Oscillator instance.
func NewStochasticOscillator[T helper.Number]() *StochasticOscillator[T] {
	return &StochasticOscillator[T]{
		Max: trend.NewMovingMaxWithPeriod[T](DefaultStochasticOscillatorMaxAndMinPeriod),
		Min: trend.NewMovingMinWithPeriod[T](DefaultStochasticOscillatorMaxAndMinPeriod),
		Sma: trend.NewSmaWithPeriod[T](DefaultStochasticOscillatorPeriod),
	}
}

// Compute function takes a channel of numbers and computes the Stochastic Oscillator. Returns k and d.
func (s *StochasticOscillator[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T) {
	//	K = (Closing - Lowest Low) / (Highest High - Lowest Low) * 100
	//	D = 3-Period SMA of K
	lowestSplice := helper.Duplicate(
		s.Min.Compute(lows),
		2,
	)

	highest := s.Max.Compute(highs)

	closings = helper.Skip(closings, s.Min.IdlePeriod())

	kSplice := helper.Duplicate(
		helper.MultiplyBy(
			helper.Divide(
				helper.Subtract(closings, lowestSplice[0]),
				helper.Subtract(highest, lowestSplice[1]),
			),
			100,
		),
		2,
	)

	d := s.Sma.Compute(kSplice[0])
	kSplice[1] = helper.Skip(kSplice[1], s.Sma.IdlePeriod())

	return kSplice[1], d
}

// IdlePeriod is the initial period that Stochastic Oscillator won't yield any results.
func (s *StochasticOscillator[T]) IdlePeriod() int {
	return s.Max.IdlePeriod() + s.Sma.IdlePeriod()
}
