// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

// WeightedClose represents the parameters for calculating the Weighted Close indicator.
//
//	Weighted Close = (High + Low + (Close * 2)) / 4
//
// Example:
//
//	weightedClose := trend.NewWeightedClose[float64]()
//	result := weightedClose.Compute(highs, lows, closes)
type WeightedClose[T helper.Number] struct {
}

// NewWeightedClose function initializes a new Weighted Close instance with the default parameters.
func NewWeightedClose[T helper.Number]() *WeightedClose[T] {
	return &WeightedClose[T]{}
}

// Compute function takes a channel of numbers and computes the Weighted Close over the specified period.
func (*WeightedClose[T]) Compute(highs, lows, closes <-chan T) <-chan T {
	return helper.Operate3(highs, lows, closes, func(high, low, close T) T {
		// Weighted Close = (High + Low + (Close * 2)) / 4
		return (high + low + (close * 2)) / 4
	})
}

// IdlePeriod is the initial period that Weighted Close yield any results.
func (*WeightedClose[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the Weighted Close.
func (*WeightedClose[T]) String() string {
	return "Weighted Close"
}
