// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Ma represents the interface for the Moving Average (MA) indicators.
type Ma[T helper.Number] interface {
	// Compute function takes a channel of numbers and computes the MA.
	Compute(<-chan T) <-chan T

	// IdlePeriod is the initial period that MA won't yield any results.
	IdlePeriod() int

	// String is the string representation of the MA instance.
	String() string
}

// MaWithContext represents the interface for the Moving Average (MA) indicators
// that support context-aware computation.
type MaWithContext[T helper.Number] interface {
	Ma[T]
	ComputeWithContext(context.Context, <-chan T) <-chan T
}

// ComputeMaWithContext computes moving average of a channel with context.
func ComputeMaWithContext[T helper.Number](ctx context.Context, ma Ma[T], c <-chan T) <-chan T {
	if mac, ok := ma.(MaWithContext[T]); ok {
		return mac.ComputeWithContext(ctx, c)
	}
	return ma.Compute(c)
}
