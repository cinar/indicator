// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultAroonPeriod is the default Aroon period of 25.
	DefaultAroonPeriod = 25
)

// Aroon represent the configuration for calculating the Aroon indicator. It is
// a technical analysis tool that gauges trend direction and strength in asset
// prices. It comprises two lines: Aroon Up and Aroon Down. Aroon Up measures
// uptrend strength, while Aroon Down measures downtrend strength. When Aroon
// Up exceeds Aroon Down, it suggests a bullish trend; when Aroon Down
// surpasses Aroon Up, it indicates a bearish trend.
//
//	Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
//	Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
//
// Example:
//
//	aroon := trend.NewAroon[float64]()
//	aroon.Period = 25
//
//	result := aroon.Compute(c)
type Aroon[T helper.Number] struct {
	// Period is the period to use.
	Period int
}

// NewAroon function initializes a new Aroon instance
// with the default parameters.
func NewAroon[T helper.Number]() *Aroon[T] {
	return &Aroon[T]{
		Period: DefaultAroonPeriod,
	}
}

// Compute function takes a channel of numbers and computes the Aroon
// over the specified period.
func (a *Aroon[T]) Compute(high, low <-chan T) (<-chan T, <-chan T) {
	max := NewMovingMax[T]()
	max.Period = a.Period

	min := NewMovingMin[T]()
	min.Period = a.Period

	sinceLastHigh := helper.Since[T, T](max.Compute(high))
	sinceLastLow := helper.Since[T, T](min.Compute(low))

	// Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
	aroonUp := helper.MultiplyBy(sinceLastHigh, -1)
	aroonUp = helper.IncrementBy(aroonUp, T(a.Period))
	aroonUp = helper.DivideBy(aroonUp, T(a.Period))
	aroonUp = helper.MultiplyBy(aroonUp, 100)
	aroonUp = helper.RoundDigits(aroonUp, 0)

	// Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
	aroonDown := helper.MultiplyBy(sinceLastLow, -1)
	aroonDown = helper.IncrementBy(aroonDown, T(a.Period))
	aroonDown = helper.DivideBy(aroonDown, T(a.Period))
	aroonDown = helper.MultiplyBy(aroonDown, 100)
	aroonDown = helper.RoundDigits(aroonDown, 0)

	return aroonUp, aroonDown
}
