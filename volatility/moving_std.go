// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"math"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultMovingStdPeriod is the default time period for Moving Standard Deviation.
	DefaultMovingStdPeriod = 1
)

// MovingStd represents the configuration parameters for calculating the Moving Standard Deviation
// over the specified period.
//
//	Std = Sqrt(1/Period * Sum(Pow(value - sma), 2))
type MovingStd[T helper.Number] struct {
	// Time period.
	Period int
}

// NewMovingStd function initializes a new Moving Standard Deviation instance with the default parameters.
func NewMovingStd[T helper.Number]() *MovingStd[T] {
	return NewMovingStdWithPeriod[T](DefaultMovingStdPeriod)
}

// NewMovingStdWithPeriod function initializes a new Moving Standard Deviation instance with the given period.
func NewMovingStdWithPeriod[T helper.Number](period int) *MovingStd[T] {
	return &MovingStd[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the Moving Standard Deviation over the specified period.
func (m *MovingStd[T]) Compute(c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	//	Std = Sqrt(1/Period * Sum(Pow(value - sma), 2))
	go func() {
		defer close(result)

		ring := helper.NewRing[T](m.Period)
		sum := T(0)

		for n := range c {
			sum -= ring.Put(n)
			sum += n

			if ring.IsFull() {
				sma := sum / T(m.Period)
				sum2 := T(0)

				for i := 0; i < m.Period; i++ {
					sum2 += T(math.Pow(float64(ring.At(i)-sma), 2))
				}

				std := T(math.Sqrt(float64(sum2 / T(m.Period))))
				result <- std
			}
		}
	}()

	return result
}

// IdlePeriod is the initial period that Moving Standard Deviation won't yield any results.
func (m *MovingStd[T]) IdlePeriod() int {
	return m.Period - 1
}
