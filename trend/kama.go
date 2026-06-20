// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

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

// ComputeWithContext function takes a channel of numbers and computes the KAMA over the specified period, supporting context cancellation.
func (k *Kama[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 3)

	//	Direction = Abs(Close - Previous Close Period Ago)
	directions := helper.Abs(
		helper.Change(closingsSplice[0], k.ErPeriod),
	)

	//	Volatility = MovingSum(Period, Abs(Close - Previous Close))
	movingSum := NewMovingSumWithPeriod[T](k.ErPeriod)
	volatilitys := movingSum.ComputeWithContext(ctx,
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
	closingsSplice[2] = helper.SkipWithContext(ctx, closingsSplice[2], k.ErPeriod-1)

	kama := make(chan T)

	go func() {
		defer close(kama)
		defer helper.DrainWithContext(ctx, scs)
		defer helper.DrainWithContext(ctx, closingsSplice[2])

		var prevKama T
		var ok bool
		select {
		case <-ctx.Done():
			return
		case prevKama, ok = <-closingsSplice[2]:
			if !ok {
				return
			}
		}

		for {
			// Note: We assume a synchronous upstream pipeline where the inputs on closingsSplice[2] and scs are aligned.
			// If cancellation happens between the two reads, the deferred DrainWithContext will consume the remaining items to avoid leaks.
			var closing T
			select {
			case <-ctx.Done():
				return
			case closing, ok = <-closingsSplice[2]:
				if !ok {
					return
				}
			}

			var sc T
			select {
			case <-ctx.Done():
				return
			case sc, ok = <-scs:
				if !ok {
					return
				}
			}

			prevKama = prevKama + sc*(closing-prevKama)
			select {
			case <-ctx.Done():
				return
			case kama <- prevKama:
			}
		}
	}()

	return kama
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (k *Kama[T]) Compute(closings <-chan T) <-chan T {
	return k.ComputeWithContext(context.Background(), closings)
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
