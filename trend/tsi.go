// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultTsiFirstSmoothingPeriod is the default first smoothing period of 25.
	DefaultTsiFirstSmoothingPeriod = 25

	// DefaultTsiSecondSmoothingPeriod is the default second smoothing period of 13.
	DefaultTsiSecondSmoothingPeriod = 13
)

// Tsi represents the parameters needed to calculate the True Strength Index (TSI). It is a technical momentum
// oscillator used in financial analysis. The TSI helps identify trends and potential trend reversals.
//
//	PCDS = Ema(13, Ema(25, (Current - Prior)))
//	APCDS = Ema(13, Ema(25, Abs(Current - Prior)))
//	TSI = (PCDS / APCDS) * 100
//
// Example:
//
//	tsi := trend.NewTsi[float64]()
//	result := tsi.Compute(closings)
type Tsi[T helper.Number] struct {
	// FirstSmoothing is the first smoothing moving average.
	FirstSmoothing Ma[T]

	// SecondSmoothing is the second smoothing moving average.
	SecondSmoothing Ma[T]
}

// NewTsi function initializes a new TSI instance with the default parameters.
func NewTsi[T helper.Number]() *Tsi[T] {
	return NewTsiWith[T](
		DefaultTsiFirstSmoothingPeriod,
		DefaultTsiSecondSmoothingPeriod,
	)
}

// NewTsiWith function initializes a new TSI instance with the given parameters.
func NewTsiWith[T helper.Number](firstSmoothingPeriod, secondSmoothingPeriod int) *Tsi[T] {
	return &Tsi[T]{
		FirstSmoothing:  NewEmaWithPeriod[T](firstSmoothingPeriod),
		SecondSmoothing: NewEmaWithPeriod[T](secondSmoothingPeriod),
	}
}

// Compute function takes a channel of numbers and computes the TSI over the specified period.
func (t *Tsi[T]) Compute(closings <-chan T) <-chan T {
	// Price change
	pcsSplice := helper.Duplicate(
		helper.Change(closings, 1),
		2,
	)

	//	PCDS = Ema(13, Ema(25, (Current - Prior)))
	pcds := t.FirstSmoothing.Compute(
		t.SecondSmoothing.Compute(
			pcsSplice[0],
		),
	)

	// APCDS = Ema(13, Ema(25, Abs(Current - Prior)))
	apcds := t.FirstSmoothing.Compute(
		t.SecondSmoothing.Compute(
			helper.Abs(pcsSplice[1]),
		),
	)

	// TSI = (PCDS / APCDS) * 100
	tsi := helper.MultiplyBy(
		helper.Divide(
			pcds,
			apcds,
		),
		T(100),
	)

	return tsi
}

// IdlePeriod is the initial period that TSI yield any results.
func (t *Tsi[T]) IdlePeriod() int {
	return t.FirstSmoothing.IdlePeriod() + t.SecondSmoothing.IdlePeriod() + 1
}

// String is the string representation of the TSI.
func (t *Tsi[T]) String() string {
	return fmt.Sprintf("TSI(%s,%s)",
		t.FirstSmoothing.String(),
		t.SecondSmoothing.String(),
	)
}
