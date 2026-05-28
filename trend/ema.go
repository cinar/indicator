// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultEmaPeriod is the default EMA period of 20.
	DefaultEmaPeriod = 20

	// DefaultEmaSmoothing is the default EMA smooting of 2.
	DefaultEmaSmoothing = 2
)

// Ema represents the parameters for calculating the Exponential Moving Average.
//
// Example:
//
//	ema := trend.NewEma[float64]()
//	ema.Period = 10
//
//	result := ema.Compute(c)
type Ema[T helper.Number] struct {
	// Time period.
	Period int

	// Smoothing constant.
	Smoothing T
}

// NewEma function initializes a new EMA instance with the default parameters.
func NewEma[T helper.Number]() *Ema[T] {
	return &Ema[T]{
		Period:    DefaultEmaPeriod,
		Smoothing: DefaultEmaSmoothing,
	}
}

// NewEmaWithPeriod function initializes a new EMA instance with the given period.
func NewEmaWithPeriod[T helper.Number](period int) *Ema[T] {
	ema := NewEma[T]()
	ema.Period = period

	return ema
}

// ComputeWithContext function takes a channel of numbers and computes the EMA over the specified period, supporting context cancellation.
func (e *Ema[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Initial EMA value is the SMA.
		sma := NewSma[T]()
		sma.Period = e.Period

		var before T
		var ok bool
		select {
		case <-ctx.Done():
			return
		case before, ok = <-sma.ComputeWithContext(ctx, helper.HeadWithContext(ctx, c, e.Period)):
			if !ok {
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case result <- before:
		}

		multiplier := e.Smoothing / T(e.Period+1)

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				before = (n-before)*multiplier + before
				select {
				case <-ctx.Done():
					return
				case result <- before:
				}
			}
		}
	}()

	return result
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (e *Ema[T]) Compute(c <-chan T) <-chan T {
	return e.ComputeWithContext(context.Background(), c)
}

// IdlePeriod is the initial period that EMA yield any results.
func (e *Ema[T]) IdlePeriod() int {
	return e.Period - 1
}

// String is the string representation of the EMA.
func (e *Ema[T]) String() string {
	return fmt.Sprintf("EMA(%d)", e.Period)
}
