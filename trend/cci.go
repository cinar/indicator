// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"
	"fmt"
	"math"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultCciPeriod is the default time period for CCI.
	DefaultCciPeriod = 20
)

// Cci represents the configuration parameters for calculating the Commodity Channel Index (CCI). CCI is a
// momentum-based oscillator used to help determine when an investment vehicle is reaching a condition of
// being overbought or oversold.
//
//	Moving Average = Sma(Period, Typical Price)
//	Mean Deviation = Sum(Abs(Typical Price[j] - Moving Average)) / Period, for j over the window
//	CCI = (Typical Price - Moving Average) / (0.015 * Mean Deviation)
//
// The mean deviation measures every typical price in the window against
// the current moving average, as in Lambert's definition and TA-Lib CCI.
// When the mean deviation is zero, CCI is 0.
//
// Example:
//
//	cci := trend.NewCci[float64]()
//	cci.Period = 20
//	values = cci.Compute(highs, lows, closings)
type Cci[T helper.Float] struct {
	// Time period.
	Period int
}

// NewCci function initializes a new CCI instance with the default parameters.
func NewCci[T helper.Float]() *Cci[T] {
	return NewCciWithPeriod[T](DefaultCciPeriod)
}

// NewCciWithPeriod function initializes a new CCI instance with the given period.
func NewCciWithPeriod[T helper.Float](period int) *Cci[T] {
	return &Cci[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the CCI.
func (c *Cci[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) <-chan T {
	typicalPrice := NewTypicalPrice[T]()
	tps := typicalPrice.ComputeWithContext(ctx, highs, lows, closings)

	ccis := helper.WindowWithContext(ctx, tps, c.cci, c.Period)

	return helper.SkipWithContext(ctx, ccis, c.IdlePeriod())
}

// cci computes the CCI of the newest value in the given window. The
// aggregation is order independent except for the newest value, which is
// located through the window's rotation offset.
func (c *Cci[T]) cci(window []T, offset int) T {
	n := T(len(window))

	var sum T
	for _, tp := range window {
		sum += tp
	}

	mean := sum / n

	var deviation T
	for _, tp := range window {
		deviation += T(math.Abs(float64(tp - mean)))
	}

	deviation /= n

	if deviation == 0 {
		return 0
	}

	newest := window[(offset+len(window)-1)%len(window)]

	return (newest - mean) / (0.015 * deviation)
}

// IdlePeriod is the initial period that CCI won't yield any results.
func (c *Cci[T]) IdlePeriod() int {
	return c.Period - 1
}

// String is the string representation of the CCI.
func (c *Cci[T]) String() string {
	return fmt.Sprintf("CCI(%d)", c.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (c *Cci[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	return c.ComputeWithContext(context.Background(), highs, lows, closings)
}
