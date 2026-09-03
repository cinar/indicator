// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

// PercentB represents the parameters for calculating the %B indicator.
//
//	%B = (Close - Lower Band) / (Upper Band - Lower Band)
//
// %B is expressed on a 0-1 scale locating price within the bands (0 at
// the lower band, 1 at the upper band). When the bands collapse to zero
// width (upper == lower, i.e. zero rolling standard deviation — a flat
// price window), price sits exactly at that single collapsed band value.
// The ratio is an undefined 0/0, so it falls back to the neutral
// midpoint 0.5.
type PercentB[T helper.Float] struct {
	// BollingerBands is the underlying Bollinger Bands indicator used for calculations.
	BollingerBands *BollingerBands[T]
}

// NewPercentB function initializes a new %B instance with the default parameters.
func NewPercentB[T helper.Float]() *PercentB[T] {
	return NewPercentBWithPeriod[T](DefaultBollingerBandsPeriod)
}

// NewPercentBWithPeriod function initializes a new %B instance with the given period.
func NewPercentBWithPeriod[T helper.Float](period int) *PercentB[T] {
	return &PercentB[T]{
		BollingerBands: NewBollingerBandsWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the %B over the specified period.
func (p *PercentB[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)

	// Compute the Bollinger Bands
	upperBands, middleBands, lowerBands := p.BollingerBands.ComputeWithContext(ctx, closingsSplice[0])

	// Skip closings until the Bollinger Bands are available
	closingsSplice[1] = helper.SkipWithContext(ctx, closingsSplice[1], p.BollingerBands.IdlePeriod())

	// Drain the middle bands
	go helper.DrainWithContext(ctx, middleBands)

	return helper.Operate3WithContext(ctx, upperBands, lowerBands, closingsSplice[1], func(upperBand, lowerBand, closing T) T {
		denom := upperBand - lowerBand
		if denom == 0 {
			return 0.5
		}

		// %B = (Close - Lower Band) / (Upper Band - Lower Band)
		return (closing - lowerBand) / denom
	})
}

// IdlePeriod is the initial period that %B yield any results.
func (p *PercentB[T]) IdlePeriod() int {
	return p.BollingerBands.IdlePeriod()
}

// String is the string representation of the %B.
func (p *PercentB[T]) String() string {
	return fmt.Sprintf("%%B(%d)", p.BollingerBands.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (p *PercentB[T]) Compute(closings <-chan T) <-chan T {
	return p.ComputeWithContext(context.Background(), closings)
}
