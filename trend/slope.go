// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultSlopePeriod is the default Slope period.
	DefaultSlopePeriod = 14
)

// Slope represents the configuration parameters for calculating the
// Rate of Change Slope indicator.
//
//	Slope = (Current Price - Price n periods ago) / n
type Slope[T helper.Number] struct {
	// Time period.
	Period int
}

// NewSlope function initializes a new Slope instance with the default parameters.
func NewSlope[T helper.Number]() *Slope[T] {
	return NewSlopeWithPeriod[T](DefaultSlopePeriod)
}

// NewSlopeWithPeriod function initializes a new Slope instance with the given parameters.
func NewSlopeWithPeriod[T helper.Number](period int) *Slope[T] {
	if period <= 0 {
		period = DefaultSlopePeriod
	}

	return &Slope[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the Slope.
func (s *Slope[T]) Compute(values <-chan T) <-chan T {
	window := helper.NewRing[T](s.Period)

	slopes := helper.Map(values, func(value T) T {
		var result T

		if window.IsFull() {
			previous, ok := window.Get()
			if ok {
				result = (value - previous) / T(s.Period)
			}
		}

		window.Put(value)

		return result
	})

	return helper.Skip(slopes, s.IdlePeriod())
}

// IdlePeriod is the initial period that Slope won't yield any results.
func (s *Slope[T]) IdlePeriod() int {
	return s.Period
}

// String is the string representation of the Slope.
func (s *Slope[T]) String() string {
	return fmt.Sprintf("SLOPE(%d)", s.Period)
}
