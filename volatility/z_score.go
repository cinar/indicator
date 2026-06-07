// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultZScorePeriod is the default period for Z-Score.
	DefaultZScorePeriod = 20
)

// ZScore represents the configuration parameters for Z-Score.
// It measures how many standard deviations price is away from its SMA.
//
//	Z-Score = (Price - SMA) / StdDev
//
// Example:
//
//	z := NewZScore[float64]()
//	z.Compute(c)
type ZScore[T helper.Number] struct {
	// Period is the time period.
	Period int
}

// NewZScore function initializes a new Z-Score instance with default parameters.
func NewZScore[T helper.Number]() *ZScore[T] {
	return NewZScoreWithPeriod[T](DefaultZScorePeriod)
}

// NewZScoreWithPeriod function initializes a new Z-Score instance with the given period.
func NewZScoreWithPeriod[T helper.Number](period int) *ZScore[T] {
	return &ZScore[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Z-Score over the specified period.
func (z *ZScore[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	cs := helper.DuplicateWithContext(ctx, c, 3)

	sma := trend.NewSmaWithPeriod[T](z.Period)
	std := NewMovingStdWithPeriod[T](z.Period)

	smaChan := sma.ComputeWithContext(ctx, cs[0])
	stdChan := std.ComputeWithContext(ctx, cs[1])
	priceChan := helper.SkipWithContext(ctx, cs[2], z.IdlePeriod())

	return helper.DivideWithContext(ctx, helper.SubtractWithContext(ctx, priceChan, smaChan), stdChan)
}

// IdlePeriod is the initial period that Z-Score won't yield any results.
func (z *ZScore[T]) IdlePeriod() int {
	return z.Period - 1
}

// String is the string representation of Z-Score.
func (z *ZScore[T]) String() string {
	return fmt.Sprintf("ZSCORE(%d)", z.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (z *ZScore[T]) Compute(c <-chan T) <-chan T {
	return z.ComputeWithContext(context.Background(), c)
}
