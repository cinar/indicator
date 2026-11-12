// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

// DefaultDpoPeriod is the default period for DPO calculation.
const DefaultDpoPeriod = 20

// Dpo computes the Detrended Price Oscillator. It compares a past price
// with the current simple moving average, which removes the trend from
// the price:
//
//	k = floor(period/2) + 1
//	DPO[t] = Price[t - k] - SMA[t]
//
// This matches TA-Lib DPO.
//
// Example:
//
//	dpo := trend.NewDpoWithPeriod[float64](20)
//	out := dpo.Compute(c)
type Dpo[T helper.Float] struct {
	// Period is the SMA window length. Typical default is 20.
	// Note: values <= 1 are considered invalid and will be replaced with DefaultDpoPeriod by constructors.
	period int
}

// NewDpo creates a new DPO instance with default parameters.
func NewDpo[T helper.Float]() *Dpo[T] {
	return &Dpo[T]{
		period: DefaultDpoPeriod,
	}
}

// NewDpoWithPeriod initializes a new DPO instance with the given period.
// Periods <= 1 are clamped to DefaultDpoPeriod.
func NewDpoWithPeriod[T helper.Float](period int) *Dpo[T] {
	if period <= 1 {
		period = DefaultDpoPeriod
	}

	return &Dpo[T]{
		period: period,
	}
}

// ComputeWithContext calculates the DPO indicator over the input price channel.
func (d *Dpo[T]) ComputeWithContext(ctx context.Context, closing <-chan T) <-chan T {
	k := d.shift()
	idle := d.IdlePeriod()
	dup := helper.DuplicateWithContext(ctx, closing, 2)

	sma := NewSmaWithPeriod[T](d.period)
	smaOut := helper.SkipWithContext(ctx, sma.ComputeWithContext(ctx, dup[0]), idle-sma.IdlePeriod())

	// Price[t - k] is aligned with SMA[t] by dropping the first idle-k
	// prices and the last k prices.
	shiftedClosing := helper.SkipLastWithContext(ctx, helper.SkipWithContext(ctx, dup[1], idle-k), k)

	return helper.SubtractWithContext(ctx, shiftedClosing, smaOut)
}

// shift returns the number of bars, k, that the price is shifted back by.
func (d *Dpo[T]) shift() int {
	return d.period/2 + 1
}

// IdlePeriod returns the number of leading samples to discard before the first DPO value is available.
func (d *Dpo[T]) IdlePeriod() int {
	return max(d.period-1, d.shift())
}

// String is the string representation of the DPO.
func (d *Dpo[T]) String() string {
	return fmt.Sprintf("DPO(%d)", d.period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (d *Dpo[T]) Compute(closing <-chan T) <-chan T {
	return d.ComputeWithContext(context.Background(), closing)
}
