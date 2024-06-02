// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultPoPeriod is the default period for the Projection Oscillator (PO).
	DefaultPoPeriod = 14
)

// Po represents the configuration parameters for calculating the Projection Oscillator (PO). It uses the linear
// regression slope, along with highs and lows. Period defines the moving window to calculates the PO.
//
//	PL = Min(period, (high + MLS(period, x, high)))
//	PH = Max(period, (low + MLS(period, x, low)))
//	PO = 100 * (Closing - PL) / (PH - PL)
//
// Example:
//
//	po := volatility.NewPo()
//	ps := po.Compute(highs, lows, closings)
type Po[T helper.Number] struct {
	// Mls is the Moving Least Square instance.
	mls *trend.Mls[T]

	// Min is the Moving Min instance.
	min *trend.MovingMin[T]

	// Max is the Moving Max instance.
	max *trend.MovingMax[T]
}

// NewPo function initializes a new PO instance with the default parameters.
func NewPo[T helper.Number]() *Po[T] {
	return NewPoWithPeriod[T](DefaultPoPeriod)
}

// NewPoWithPeriod function initializes a new PO instance with the given period.
func NewPoWithPeriod[T helper.Number](period int) *Po[T] {
	return &Po[T]{
		mls: trend.NewMlsWithPeriod[T](period),
		min: trend.NewMovingMinWithPeriod[T](period),
		max: trend.NewMovingMaxWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the PO over the specified period.
func (p *Po[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	highsSplice := helper.Duplicate(highs, 2)
	lowsSplice := helper.Duplicate(lows, 2)
	closingsSplice := helper.Duplicate(closings, 2)

	xSplice := helper.Duplicate(
		helper.Count(T(1), closingsSplice[0]),
		2,
	)

	// PL = Min(period, (high + MLS(period, x, high)))
	plM, plB := p.mls.Compute(xSplice[0], highsSplice[0])
	go helper.Drain(plB)

	highsSplice[1] = helper.Skip(highsSplice[1], p.mls.IdlePeriod())

	plSplice := helper.Duplicate(
		p.min.Compute(
			helper.Add(
				highsSplice[1],
				plM,
			),
		),
		2,
	)

	// PH = Max(period, (low + MLS(period, x, low)))
	phM, phB := p.mls.Compute(xSplice[1], lowsSplice[0])
	go helper.Drain(phB)

	lowsSplice[1] = helper.Skip(lowsSplice[1], p.mls.IdlePeriod())

	ph := p.max.Compute(
		helper.Add(
			lowsSplice[1],
			phM,
		),
	)

	// PO = 100 * (Closing - PL) / (PH - PL)
	closingsSplice[1] = helper.Skip(closingsSplice[1], p.mls.IdlePeriod()+p.min.IdlePeriod())

	po := helper.MultiplyBy(
		helper.Divide(
			helper.Subtract(
				closingsSplice[1],
				plSplice[0],
			),
			helper.Subtract(
				ph,
				plSplice[1],
			),
		),
		T(100),
	)

	return po
}

// IdlePeriod is the initial period that PO won't yield any results.
func (p *Po[T]) IdlePeriod() int {
	return p.mls.IdlePeriod() + p.min.IdlePeriod()
}
