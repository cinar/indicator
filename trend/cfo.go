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
	// DefaultCfoPeriod is the default CFO period of 14.
	DefaultCfoPeriod = 14
)

// Cfo represents the configuration parameters for calculating the
// Chande Forecast Oscillator (CFO). CFO is a momentum indicator
// that measures the difference between a security's price and its
// linear regression forecast.
//
//	CFO = ((Price - Forecast) / Price) * 100
//
// Example:
//
//	cfo := trend.NewCfo[float64]()
//	result := cfo.Compute(c)
type Cfo[T helper.Float] struct {
	// Mlr is the Moving Linear Regression instance.
	Mlr *Mlr[T]
}

// NewCfo function initializes a new CFO instance with the default parameters.
func NewCfo[T helper.Float]() *Cfo[T] {
	return NewCfoWithPeriod[T](DefaultCfoPeriod)
}

// NewCfoWithPeriod function initializes a new CFO instance with the given period.
func NewCfoWithPeriod[T helper.Float](period int) *Cfo[T] {
	return &Cfo[T]{
		Mlr: NewMlrWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the CFO.
func (c *Cfo[T]) ComputeWithContext(ctx context.Context, closing <-chan T) <-chan T {
	closingSplices := helper.DuplicateWithContext(ctx, closing, 3)

	x := helper.CountWithContext(ctx, T(0), closingSplices[0])
	forecast := c.Mlr.ComputeWithContext(ctx, x, closingSplices[1])

	closingPriceSplice := helper.DuplicateWithContext(ctx, helper.SkipWithContext(ctx, closingSplices[2], c.IdlePeriod()), 2)

	return helper.MultiplyByWithContext(ctx, helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, closingPriceSplice[0],
		forecast,
	),
		closingPriceSplice[1],
	),
		T(100),
	)
}

// IdlePeriod is the initial period that CFO won't yield any results.
func (c *Cfo[T]) IdlePeriod() int {
	return c.Mlr.IdlePeriod()
}

// String is the string representation of the CFO.
func (c *Cfo[T]) String() string {
	return fmt.Sprintf("CFO(%d)", c.Mlr.Mls.Sum.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (c *Cfo[T]) Compute(closing <-chan T) <-chan T {
	return c.ComputeWithContext(context.Background(), closing)
}
