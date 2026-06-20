// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultQstickPeriod is the default period for the Qstick SMA.
	DefaultQstickPeriod = 20
)

// Qstick represents the configuration parameter for calculating the
// Qstick indicator. Qstick is a momentum indicator used to identify
// an asset's trend by looking at the SMA of the difference between
// its closing and opening.
//
// A Qstick above zero indicates increasing buying pressure, while
// a Qstick below zero indicates increasing selling pressure.
//
//	QS = SMA(Closings - Openings)
//
// Example:
//
//	qstick := momentum.Qstick[float64]()
//	qstick.Sma.Period = 50
//
//	values := qstick.Compute(openings, closings)
type Qstick[T helper.Number] struct {
	Sma *trend.Sma[T]
}

// NewQstick function initializes a new QStick instance.
func NewQstick[T helper.Number]() *Qstick[T] {
	qstick := &Qstick[T]{
		Sma: trend.NewSma[T](),
	}

	qstick.Sma.Period = DefaultQstickPeriod

	return qstick
}

// ComputeWithContext function takes a channel of numbers and computes the Qstick.
func (q *Qstick[T]) ComputeWithContext(ctx context.Context, openings, closings <-chan T) <-chan T {
	qstick := helper.SubtractWithContext(ctx, closings, openings)
	qstick = q.Sma.ComputeWithContext(ctx, qstick)

	return qstick
}

// IdlePeriod is the initial period that Qstick won't yield any results.
func (q *Qstick[T]) IdlePeriod() int {
	return q.Sma.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (q *Qstick[T]) Compute(openings, closings <-chan T) <-chan T {
	return q.ComputeWithContext(context.Background(), openings, closings)
}
