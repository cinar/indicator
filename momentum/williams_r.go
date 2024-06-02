// WilliamsRpyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultWilliamsRPeriod is the default period for the Williams R.
	DefaultWilliamsRPeriod = 14
)

// WilliamsR represents the configuration parameter for calculating the Williams %R, or just %R. It is a technical
// analysis oscillator showing the current closing price in relation to the high and low of the past N days (for a
// given N). It was developed by a publisher and promoter of trading materials, Larry Williams. Its purpose is to
// tell whether a stock or commodity market is trading near the high or the low, or somewhere in between, of
// its recent trading range. Buy when -80 and below. Sell when -20 and above.
//
//	WR = (Highest High - Closing) / (Highest High - Lowest Low) * -100.
//
// Example:
//
//	wr := momentum.WilliamsR[float64]()
//	values := wr.Compute(highs, lows, closings)
type WilliamsR[T helper.Number] struct {
	// Max is the Moving Max instance.
	Max *trend.MovingMax[T]

	// Min is the Moving Min instance.
	Min *trend.MovingMin[T]
}

// NewWilliamsR function initializes a new Williams R instance.
func NewWilliamsR[T helper.Number]() *WilliamsR[T] {
	return &WilliamsR[T]{
		Max: trend.NewMovingMaxWithPeriod[T](DefaultWilliamsRPeriod),
		Min: trend.NewMovingMinWithPeriod[T](DefaultWilliamsRPeriod),
	}
}

// Compute function takes a channel of numbers and computes the Williams R.
func (w *WilliamsR[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	highestSplice := helper.Duplicate(
		w.Max.Compute(highs),
		2,
	)

	lowest := w.Min.Compute(lows)

	closings = helper.Skip(closings, w.Max.IdlePeriod())

	return helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(highestSplice[0], closings),
			helper.Subtract(highestSplice[1], lowest),
		),
		-100,
	)
}

// IdlePeriod is the initial period that Williams R won't yield any results.
func (w *WilliamsR[T]) IdlePeriod() int {
	return w.Max.IdlePeriod()
}
