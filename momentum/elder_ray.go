// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultElderRayPeriod is the default period for Elder-Ray Index.
	DefaultElderRayPeriod = 13
)

// ElderRay represents the configuration parameters for calculating the Elder-Ray Index.
// Developed by Alexander Elder, the Elder-Ray Index measures buying and selling pressure
// in the market. It consists of two separate indicators: Bull Power and Bear Power.
//
//	Bull Power = High - n-period EMA
//	Bear Power = Low - n-period EMA
//
// Example:
//
//	er := momentum.NewElderRay[float64]()
//	bullPower, bearPower := er.Compute(highs, lows, closings)
type ElderRay[T helper.Number] struct {
	// Period is the time period.
	Period int
}

// NewElderRay function initializes a new Elder-Ray Index instance with the default parameters.
func NewElderRay[T helper.Number]() *ElderRay[T] {
	return NewElderRayWithPeriod[T](DefaultElderRayPeriod)
}

// NewElderRayWithPeriod function initializes a new Elder-Ray Index instance with the given period.
func NewElderRayWithPeriod[T helper.Number](period int) *ElderRay[T] {
	return &ElderRay[T]{
		Period: period,
	}
}

// Compute function takes channels of highs, lows, and closings and computes the Elder-Ray Index.
// Returns bullPower and bearPower channels.
func (e *ElderRay[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T) {
	ema := trend.NewEmaWithPeriod[T](e.Period)
	emas := helper.Duplicate(ema.Compute(closings), 2)

	highs = helper.Skip(highs, e.IdlePeriod())
	lows = helper.Skip(lows, e.IdlePeriod())

	bullPower := helper.Subtract(highs, emas[0])
	bearPower := helper.Subtract(lows, emas[1])

	return bullPower, bearPower
}

// IdlePeriod is the initial period that Elder-Ray Index won't yield any results.
func (e *ElderRay[T]) IdlePeriod() int {
	return e.Period - 1
}

// String is the string representation of the Elder-Ray Index.
func (e *ElderRay[T]) String() string {
	return fmt.Sprintf("Elder-Ray Index(%d)", e.Period)
}
