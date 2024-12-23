// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultSmmaPeriod is the default SMMA period of 7.
	DefaultSmmaPeriod = 7
)

// Smma represents the parameters for calculating the Smoothed Moving Average (SMMA).
//
//	SMMA[0] = SMA(N)
//	SMMA[i] = ((SMMA[i-1] * (N - 1)) + Close[i]) / N
//
// Example:
//
//	smma := trend.NewSmma[float64]()
//	smma.Period = 10
//
//	result := smma.Compute(c)
type Smma[T helper.Number] struct {
	// Time period.
	Period int
}

// NewSmma function initializes a new SMMA instance with the default parameters.
func NewSmma[T helper.Number]() *Smma[T] {
	return &Smma[T]{
		Period: DefaultSmmaPeriod,
	}
}

// NewSmmaWithPeriod function initializes a new SMMA instance with the given period.
func NewSmmaWithPeriod[T helper.Number](period int) *Smma[T] {
	smma := NewSmma[T]()
	smma.Period = period

	return smma
}

// Compute function takes a channel of numbers and computes the SMMA over the specified period.
func (s *Smma[T]) Compute(c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Initial SMMA value is the SMA.
		sma := NewSmaWithPeriod[T](s.Period)

		before := <-sma.Compute(helper.Head(c, s.Period))
		result <- before

		for n := range c {
			before = ((before*T(s.Period) - 1) + n) / T(s.Period)
			result <- before
		}
	}()

	return result
}

// IdlePeriod is the initial period that SMMA yield any results.
func (s *Smma[T]) IdlePeriod() int {
	return s.Period - 1
}

// String is the string representation of the SMMA.
func (s *Smma[T]) String() string {
	return fmt.Sprintf("SMMA(%d)", s.Period)
}
