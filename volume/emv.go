// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultEmvPeriod is the default period for the EMV.
	DefaultEmvPeriod = 14
)

// Emv holds configuration parameters for calculating the Ease of Movement (EMV). It is a volume based oscillator
// measuring the ease of price movement.
//
//	Distance Moved = ((High + Low) / 2) - ((Priod High + Prior Low) /2)
//	Box Ratio = ((Volume / 100000000) / (High - Low))
//	EMV(1) = Distance Moved / Box Ratio
//	EMV(14) = SMA(14, EMV(1))
//
// Example:
//
//	emv := volume.NewEmv[float64]()
//	result := emv.Compute(highs, lows, volumes)
type Emv[T helper.Number] struct {
	// Sma is the SMA instance.
	Sma *trend.Sma[T]
}

// NewEmv function initializes a new EMV instance with the default parameters.
func NewEmv[T helper.Number]() *Emv[T] {
	return NewEmvWithPeriod[T](DefaultEmvPeriod)
}

// NewEmvWithPeriod function initializes a new EMV instance with the given period.
func NewEmvWithPeriod[T helper.Number](period int) *Emv[T] {
	return &Emv[T]{
		Sma: trend.NewSmaWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the EMV.
func (e *Emv[T]) Compute(highs, lows, volumes <-chan T) <-chan T {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)

	//	Distance Moved = ((High + Low) / 2) - ((Priod High + Prior Low) /2)
	distanceMoved := helper.Change(
		helper.DivideBy(
			helper.Add(
				highsSplice[0],
				lowsSplice[0],
			),
			2,
		),
		1,
	)

	d := 100000000

	// Box Ratio = ((Volume / 100000000) / (High - Low))
	boxRatio := helper.Divide(
		helper.DivideBy(
			volumes,
			T(d),
		),
		helper.Subtract(
			highsSplice[1],
			lowsSplice[1],
		),
	)

	// EMV(1) = Distance Moved / Box Ratio
	// EMV(14) = SMA(14, EMV(1))
	return e.Sma.Compute(
		helper.Divide(
			distanceMoved,
			boxRatio,
		),
	)
}

// IdlePeriod is the initial period that EMV won't yield any results.
func (e *Emv[T]) IdlePeriod() int {
	return e.Sma.IdlePeriod() + 1
}
