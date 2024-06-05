// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultKamaErPeriod is the default Efficiency Ratio (ER) period of 10.
	DefaultKamaErPeriod = 10

	// DefaultKamaFastScPeriod is the default Fast Smoothing Constant (SC) period of 2.
	DefaultKamaFastScPeriod = 2

	// DefaultKamaSlowScPeriod is the default Slow Smoothing Constant (SC) period of 30.
	DefaultKamaSlowScPeriod = 30
)

// Kama represents the parameters for calculating the Kaufman's Adaptive Moving Average (KAMA).
// It is a type of moving average that adapts to market noise or volatility. It tracks prices
// closely during periods of small price swings and low noise.
//
//	Direction = Abs(Close - Previous Close Period Ago)
//	Volatility = MovingSum(Period, Abs(Close - Prior Close))
//	Efficiency Ratio (ER) = Direction / Volatility
//	Smoothing Constant (SC) = (ER * (2/(Fast + 1) - 2/(Slow + 1)) + (2/(Slow + 1)))^2
//	KAMA = Prior KAMA + SC * (Price - Prior KAMA)
//
// Example:
//
//	kama := trend.NewKama[float64]()
//	result := kama.Compute(c)
type Kama[T helper.Number] struct {
	// ErPeriod is the Efficiency Ratio time period.
	ErPeriod int

	// FastScPeriod is the Fast Smoothing Constant time period.
	FastScPeriod int

	// SlowScPeriod is the Slow Smoothing Constant time period.
	SlowScPeriod int
}

// NewKama function initializes a new KAMA instance with the default parameters.
func NewKama[T helper.Number]() *Kama[T] {
	return &Kama[T]{
		ErPeriod:     DefaultKamaErPeriod,
		FastScPeriod: DefaultKamaFastScPeriod,
		SlowScPeriod: DefaultKamaSlowScPeriod,
	}
}

// Compute function takes a channel of numbers and computes the KAMA over the specified period.
func (k *Kama[T]) Compute(c <-chan T) <-chan T {
	closingsSplice := helper.Duplicate(c, 1)

	//	Direction = Abs(Close - Previous Close Period Ago)
	directions := helper.Abs(
		helper.Change(closingsSplice[0], k.ErPeriod),
	)

	/*
		//	Volatility = MovingSum(Period, Abs(Close - Prior Close))
		movingSum := NewMovingSumWithPeriod[T](k.ErPeriod)
		volatilitys := movingSum.Compute(
			helper.Abs(
				helper.Change(closingsSplice[1], 1),
			),
		)

		//	Efficiency Ratio (ER) = Direction / Volatility
		ers := helper.Divide(directions, volatilitys)
	*/

	return directions
}

// IdlePeriod is the initial period that KAMA yield any results.
func (k *Kama[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the KAMA.
func (k *Kama[T]) String() string {
	return fmt.Sprintf("KAMA(%d,%d,%d)",
		k.ErPeriod,
		k.FastScPeriod,
		k.SlowScPeriod,
	)
}
