// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultTrimaPeriod is the default period for TRIMA.
	DefaultTrimaPeriod = 15
)

// Trima represents the configuration parameters for calculating the
// Triangular Moving Average (TRIMA).
//
// If period is even:
//
//	TRIMA = SMA(period / 2, SMA((period / 2) + 1, values))
//
// If period is odd:
//
//	TRIMA = SMA((period + 1) / 2, SMA((period + 1) / 2, values))
type Trima[T helper.Float] struct {
	// Time period.
	Period int
}

// NewTrima function initializes a new TRIMA instance
// with the default parameters.
func NewTrima[T helper.Float]() *Trima[T] {
	return NewTrimaWithPeriod[T](DefaultTrimaPeriod)
}

// NewTrimaWithPeriod function initializes a new TRIMA instance
// with the given period.
func NewTrimaWithPeriod[T helper.Float](period int) *Trima[T] {
	return &Trima[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the TRIMA
// and the signal line.
func (t *Trima[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	period1, period2 := t.calculatePeriods()

	sma1 := NewSma[T]()
	sma1.Period = period1

	sma2 := NewSma[T]()
	sma2.Period = period2

	trima := sma1.ComputeWithContext(ctx, sma2.ComputeWithContext(ctx, c))

	return trima
}

// IdlePeriod is the initial period that TRIMA won't yield any results.
func (t *Trima[T]) IdlePeriod() int {
	period1, period2 := t.calculatePeriods()
	return period1 + period2 - 2
}

// String is the string representation of the TRIMA.
func (t *Trima[T]) String() string {
	return fmt.Sprintf("TRIMA(%d)", t.Period)
}

// calculatePeriods calculates the individual periods to use based on the
// TRIMA period.
func (t *Trima[T]) calculatePeriods() (int, int) {
	var period1, period2 int

	if t.Period%2 == 0 {
		period1 = t.Period / 2
		period2 = period1 + 1
	} else {
		period1 = (t.Period + 1) / 2
		period2 = period1
	}

	return period1, period2
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (t *Trima[T]) Compute(c <-chan T) <-chan T { return t.ComputeWithContext(context.Background(), c) }
