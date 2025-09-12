// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

// DefaultDpoPeriod is the default period for DPO calculation.
const DefaultDpoPeriod = 20

// Dpo computes the Detrended Price Oscillator.
// Formula (common approximation):
// Let k = floor(period/2) + 1.
// For time index t >= period-1+k:
//
//	DPO[t] = Price[t] - SMA[t - k]
//
// Example:
//
//	dpo := trend.NewDpoWithPeriod[float64](20)
//	out := dpo.Compute(c)
type Dpo[T helper.Float] struct {
	// period is the SMA window length. Must be >= 1. Typical default is 20.
	period int
}

// NewDpo creates a new DPO instance with default parameters.
func NewDpo[T helper.Float]() *Dpo[T] {
	return &Dpo[T]{
		period: DefaultDpoPeriod,
	}
}

// NewDpoWithPeriod function initializes a new DPO instance with the given period.
func NewDpoWithPeriod[T helper.Float](period int) *Dpo[T] {
	if period <= 1 {
		period = DefaultDpoPeriod
	}

	return &Dpo[T]{
		period: period,
	}
}

// Compute calculates the DPO indicator over the input price channel.
func (d *Dpo[T]) Compute(closing <-chan T) <-chan T {
	k := d.period/2 + 1
	dup := helper.Duplicate(closing, 2)

	// compute SMA on the first duplicated stream
	sma := NewSma[T]()
	sma.Period = d.period
	smaOut := sma.Compute(dup[0])

	// align the original price stream and the SMA stream according to DPO formula
	skippedClosing := helper.Skip(dup[1], d.IdlePeriod())
	smaDelayed := helper.SkipLast(smaOut, k)

	// DPO = Price - shifted SMA
	return helper.Operate(skippedClosing, smaDelayed, func(price, shiftedSma T) T {
		return price - shiftedSma
	})
}

// IdlePeriod returns the number of leading samples to discard before the first DPO value is available.
func (d *Dpo[T]) IdlePeriod() int {
	return (d.period - 1) + (d.period/2 + 1)
}

// String is the string representation of the DPO.
func (d *Dpo[T]) String() string {
	return fmt.Sprintf("DPO(%d)", d.period)
}
