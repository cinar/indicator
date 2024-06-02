// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultKdjMinMaxPeriod is the default period for moving min
	// of low, and moving max of high.
	DefaultKdjMinMaxPeriod = 9

	// DefaultKdjSma1Period is the default period for SMA of RSV.
	DefaultKdjSma1Period = 3

	// DefaultKdjSma2Period is the default period for SMA of K.
	DefaultKdjSma2Period = 3
)

// Kdj represents the configuration parameters for calculating the
// KDJ, also known as the Random Index.  KDJ is calculated similar
// to the Stochastic Oscillator with the difference of having the
// J line. It is used to analyze the trend and entry points.
//
// The K and D lines show if the asset is overbought when they
// crosses above 80%, and oversold when they crosses below
// 20%. The J line represents the divergence.
//
//	RSV = ((Closing - Min(Low, rPeriod))
//	/ (Max(High, rPeriod) - Min(Low, rPeriod))) * 100
//
//	K = Sma(RSV, kPeriod)
//	D = Sma(K, dPeriod)
//	J = (3 * K) - (2 * D)
//
// Example:
//
//	kdj := NewKdj[float64]()
//	values := kdj.Compute(highs, lows, closings)
type Kdj[T helper.Number] struct {
	// MovingMax is the highest high.
	MovingMax *MovingMax[T]

	// MovingMin is the lowest low.
	MovingMin *MovingMin[T]

	// Sma1 is the SMA of RSV.
	Sma1 *Sma[T]

	// Sma2 is the SMA of K.
	Sma2 *Sma[T]
}

// NewKdj function initializes a new Kdj instance with the default parameters
func NewKdj[T helper.Number]() *Kdj[T] {
	kdj := &Kdj[T]{
		MovingMax: NewMovingMax[T](),
		MovingMin: NewMovingMin[T](),
		Sma1:      NewSma[T](),
		Sma2:      NewSma[T](),
	}

	kdj.MovingMax.Period = DefaultKdjMinMaxPeriod
	kdj.MovingMin.Period = DefaultKdjMinMaxPeriod
	kdj.Sma1.Period = DefaultKdjSma1Period
	kdj.Sma2.Period = DefaultKdjSma2Period

	return kdj
}

// Compute function takes a channel of numbers and computes the KDJ
// over the specified period. Returns K, D, J.
func (kdj *Kdj[T]) Compute(high, low, closing <-chan T) (<-chan T, <-chan T, <-chan T) {
	highest := kdj.MovingMax.Compute(high)
	lowests := helper.Duplicate(
		kdj.MovingMin.Compute(low),
		2,
	)

	closing = helper.Skip(closing, kdj.MovingMax.Period-1)

	rsv := helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(closing, lowests[0]),
			helper.Subtract(highest, lowests[1]),
		),
		100,
	)

	ks := helper.Duplicate(
		kdj.Sma1.Compute(rsv),
		3,
	)

	ds := helper.Duplicate(
		kdj.Sma2.Compute(ks[0]),
		2,
	)

	ks[1] = helper.Skip(ks[1], kdj.Sma2.Period-1)
	ks[2] = helper.Skip(ks[2], kdj.Sma2.Period-1)

	j := helper.Subtract(
		helper.MultiplyBy(ks[1], 3),
		helper.MultiplyBy(ds[0], 2),
	)

	return ks[2], ds[1], j
}

// IdlePeriod is the initial period that KDJ won't yield any results.
func (kdj *Kdj[T]) IdlePeriod() int {
	return kdj.MovingMax.Period + kdj.Sma1.Period + kdj.Sma2.Period - 3
}
