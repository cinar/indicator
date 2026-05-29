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
	// DefaultSmmaPeriod is the default SMMA period of 7.
	DefaultSmmaPeriod = 7
)

// Smma represents the parameters for calculating the Smoothed Moving Average (SMMA).
//
//	SMMA[0] = SMA(N)
//	SMMA[i] = ((SMMA[i-1] * (N - 1)) + Close[i]) / N
//
// Example:
//
//	smma := trend.NewSmma[float64]()
//	smma.Period = 10
//
//	result := smma.Compute(c)
type Smma[T helper.Number] struct {
	// Time period.
	Period int
}

// NewSmma function initializes a new SMMA instance with the default parameters.
func NewSmma[T helper.Number]() *Smma[T] {
	return NewSmmaWithPeriod[T](DefaultSmmaPeriod)
}

// NewSmmaWithPeriod function initializes a new SMMA instance with the given period.
func NewSmmaWithPeriod[T helper.Number](period int) *Smma[T] {
	return &Smma[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the SMMA over the specified period, supporting context cancellation.
func (s *Smma[T]) ComputeWithContext(ctx context.Context, c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Initial SMMA value is the SMA.
		sma := NewSmaWithPeriod[T](s.Period)

		var before T
		var ok bool
		select {
		case <-ctx.Done():
			return
		case before, ok = <-sma.ComputeWithContext(ctx, helper.HeadWithContext(ctx, c, s.Period)):
			if !ok {
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case result <- before:
		}

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				before = ((before * (T(s.Period) - 1)) + n) / T(s.Period)
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
func (s *Smma[T]) Compute(c <-chan T) <-chan T {
	return s.ComputeWithContext(context.Background(), c)
}

// IdlePeriod is the initial period that SMMA yield any results.
func (s *Smma[T]) IdlePeriod() int {
	return s.Period - 1
}

// String is the string representation of the SMMA.
func (s *Smma[T]) String() string {
	return fmt.Sprintf("SMMA(%d)", s.Period)
}
