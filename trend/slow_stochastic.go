// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultSlowStochasticPeriod is the default period for the Slow Stochastic indicator.
	DefaultSlowStochasticPeriod = 10

	// DefaultSlowStochasticKPeriod is the default period for the Fast %K SMA smoothing.
	DefaultSlowStochasticKPeriod = 3

	// DefaultSlowStochasticDPeriod is the default period for the Slow %D SMA.
	DefaultSlowStochasticDPeriod = 3
)

// SlowStochastic represents the configuration parameters for calculating
// the Slow Stochastic indicator. This applies additional smoothing to the
// Fast Stochastic values.
//
//	Fast %K = Stochastic(values, period)
//	Slow %K = SMA(Fast %K, kPeriod)
//	Slow %D = SMA(Slow %K, dPeriod)
//
// Example:
//
//	s := trend.NewSlowStochastic[float64]()
//	k, d := s.Compute(values)
type SlowStochastic[T helper.Number] struct {
	// Period is the period for the min/max calculation.
	Period int

	// KPeriod is the period for smoothing Fast %K to get Slow %K.
	KPeriod int

	// DPeriod is the period for smoothing Slow %K to get Slow %D.
	DPeriod int
}

// NewSlowStochastic function initializes a new SlowStochastic instance with the default parameters.
func NewSlowStochastic[T helper.Number]() *SlowStochastic[T] {
	return &SlowStochastic[T]{
		Period:  DefaultSlowStochasticPeriod,
		KPeriod: DefaultSlowStochasticKPeriod,
		DPeriod: DefaultSlowStochasticDPeriod,
	}
}

// NewSlowStochasticWithPeriod function initializes a new SlowStochastic instance with the given periods.
func NewSlowStochasticWithPeriod[T helper.Number](period, kPeriod, dPeriod int) *SlowStochastic[T] {
	return &SlowStochastic[T]{
		Period:  period,
		KPeriod: kPeriod,
		DPeriod: dPeriod,
	}
}

// Compute function takes a channel of numbers and computes the Slow Stochastic indicator.
// Returns Slow %K and Slow %D.
func (s *SlowStochastic[T]) Compute(values <-chan T) (<-chan T, <-chan T) {
	movingMin := NewMovingMinWithPeriod[T](s.Period)
	movingMax := NewMovingMaxWithPeriod[T](s.Period)

	values = helper.Buffered(values, s.Period)
	inputs := helper.Duplicate(values, 3)

	lowestSplice := helper.Duplicate(
		movingMin.Compute(inputs[0]),
		2,
	)

	highest := movingMax.Compute(inputs[1])

	skipped := helper.Skip(inputs[2], movingMin.IdlePeriod())

	fastK := helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(skipped, lowestSplice[0]),
			helper.Subtract(highest, lowestSplice[1]),
		),
		100,
	)

	slowKSma := NewSmaWithPeriod[T](s.KPeriod)
	slowK := slowKSma.Compute(fastK)

	slowKDuplicate := helper.Duplicate(slowK, 2)

	slowDSma := NewSmaWithPeriod[T](s.DPeriod)
	slowD := slowDSma.Compute(slowKDuplicate[0])

	slowKDuplicate[1] = helper.Skip(slowKDuplicate[1], s.DPeriod-1)

	return slowKDuplicate[1], slowD
}

// IdlePeriod is the initial period that Slow Stochastic won't yield any results.
func (s *SlowStochastic[T]) IdlePeriod() int {
	return s.Period + s.KPeriod + s.DPeriod - 3
}
