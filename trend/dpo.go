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
//
// Formula (common approximation):
//
//	k = period/2 + 1
//	DPO = Price - SMA(Price shifted by k)
//
// Example:
//
//	dpo := trend.NewDpo[float64]()
//	dpo.Period = 20
//	out := dpo.Compute(c)
type Dpo[T helper.Number] struct {
	Period int
}

// NewDpo creates a new DPO instance with default parameters.
func NewDpo[T helper.Number]() *Dpo[T] {
	return &Dpo[T]{
		Period: DefaultDpoPeriod,
	}
}

// NewDpoWithPeriod function initializes a new DPO instance with the given period.
func NewDpoWithPeriod[T helper.Number](period int) *Dpo[T] {
	if period < 1 {
		period = DefaultDpoPeriod
	}

	return &Dpo[T]{
		Period: period,
	}
}

// Compute calculates the DPO indicator over the input price channel.
func (d *Dpo[T]) Compute(closing <-chan T) <-chan T {
	closingSplice := helper.Duplicate(closing, 2)

	// compute SMA on the first duplicated stream
	sma := NewSma[T]()
	sma.Period = d.Period
	smaOut := sma.Compute(closingSplice[0])

	// align the original price stream and the SMA stream according to DPO formula
	skippedClosing := helper.Skip(closingSplice[1], d.IdlePeriod())
	smaDelayed := helper.SkipLast(smaOut, d.Period/2+1)

	// DPO = Price - shifted SMA
	return helper.Operate(skippedClosing, smaDelayed, func(price, shiftedSma T) T {
		return price - shiftedSma
	})
}

// IdlePeriod is the initial period that DPO yield any results.
func (d *Dpo[T]) IdlePeriod() int {
	return (d.Period - 1) + (d.Period/2 + 1)
}

// String is the string representation of the DPO.
func (d *Dpo[T]) String() string {
	return fmt.Sprintf("DPO(%d)", d.Period)
}
