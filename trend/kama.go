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
//	Volatility = MovingSum(Period, Abs(Close - Previous Close))
//	Efficiency Ratio (ER) = Direction / Volatility
//	Smoothing Constant (SC) = (ER * (2/(Fast + 1) - 2/(Slow + 1)) + (2/(Slow + 1)))^2
//	KAMA = Previous KAMA + SC * (Price - Previous KAMA)
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
	return NewKamaWith[T](
		DefaultKamaErPeriod,
		DefaultKamaFastScPeriod,
		DefaultKamaSlowScPeriod,
	)
}

// NewKamaWith function initializes a new KAMA instance with the given parameters.
func NewKamaWith[T helper.Number](erPeriod, fastScPeriod, slowScPeriod int) *Kama[T] {
	return &Kama[T]{
		ErPeriod:     erPeriod,
		FastScPeriod: fastScPeriod,
		SlowScPeriod: slowScPeriod,
	}
}

// Compute function takes a channel of numbers and computes the KAMA over the specified period.
func (k *Kama[T]) Compute(closings <-chan T) <-chan T {
	closingsSplice := helper.Duplicate(closings, 3)

	//	Direction = Abs(Close - Previous Close Period Ago)
	directions := helper.Abs(
		helper.Change(closingsSplice[0], k.ErPeriod),
	)

	//	Volatility = MovingSum(Period, Abs(Close - Previous Close))
	movingSum := NewMovingSumWithPeriod[T](k.ErPeriod)
	volatilitys := movingSum.Compute(
		helper.Abs(
			helper.Change(closingsSplice[1], 1),
		),
	)

	//	Efficiency Ratio (ER) = Direction / Volatility
	ers := helper.Divide(directions, volatilitys)

	//	Smoothing Constant (SC) = (ER * (2/(Slow + 1) - 2/(Fast + 1)) + (2/(Slow + 1)))^2
	fastSc := T(2.0) / T(k.FastScPeriod+1)
	slowSc := T(2.0) / T(k.SlowScPeriod+1)

	scs := helper.Pow(
		helper.IncrementBy(
			helper.MultiplyBy(
				ers,
				fastSc-slowSc,
			),
			slowSc,
		),
		2,
	)

	//	KAMA = Previous KAMA + SC * (Price - Previous KAMA)
	closingsSplice[2] = helper.Skip(closingsSplice[2], k.ErPeriod-1)

	kama := make(chan T)

	go func() {
		defer close(kama)
		defer helper.Drain(scs)
		defer helper.Drain(closingsSplice[2])

		prevKama, ok := <-closingsSplice[2]
		if !ok {
			return
		}

		for {
			closing, ok := <-closingsSplice[2]
			if !ok {
				break
			}

			sc := <-scs

			prevKama = prevKama + sc*(closing-prevKama)
			kama <- prevKama
		}
	}()

	return kama
}

// IdlePeriod is the initial period that KAMA yield any results.
func (k *Kama[T]) IdlePeriod() int {
	return k.ErPeriod
}

// String is the string representation of the KAMA.
func (k *Kama[T]) String() string {
	return fmt.Sprintf("KAMA(%d,%d,%d)",
		k.ErPeriod,
		k.FastScPeriod,
		k.SlowScPeriod,
	)
}
