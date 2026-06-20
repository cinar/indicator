// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Mlr represents the configuration parameters for calculating the Moving Linear Regression.
//
//	y = mx + b
//
// Example:
//
//	mlr := trend.NewMlrWithPeriod[float64](14)
//	rs := mlr.Compute(x , y)
type Mlr[T helper.Number] struct {
	// Mls is the Moving Least Square instance.
	Mls *Mls[T]
}

// NewMlrWithPeriod function initializes a new MLR instance with the given period.
func NewMlrWithPeriod[T helper.Number](period int) *Mlr[T] {
	return &Mlr[T]{
		Mls: NewMlsWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the MLR r.
func (m *Mlr[T]) ComputeWithContext(ctx context.Context, x, y <-chan T) <-chan T {
	xSplice := helper.DuplicateWithContext(ctx, x, 2)

	ms, bs := m.Mls.ComputeWithContext(ctx, xSplice[0], y)

	xSplice[1] = helper.SkipWithContext(ctx, xSplice[1], m.Mls.IdlePeriod())

	r := helper.AddWithContext(ctx, helper.MultiplyWithContext(ctx, ms,
		xSplice[1],
	),
		bs,
	)

	return r
}

// IdlePeriod is the initial period that MLR won't yield any results.
func (m *Mlr[T]) IdlePeriod() int {
	return m.Mls.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *Mlr[T]) Compute(x, y <-chan T) <-chan T {
	return m.ComputeWithContext(context.Background(), x, y)
}
