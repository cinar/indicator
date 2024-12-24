// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

// PercentB represents the parameters for calculating the %B indicator.
//
//	%B = (Close - Lower Band) / (Upper Band - Lower Band)
type PercentB[T helper.Number] struct {
	// BollingerBands is the underlying Bollinger Bands indicator used for calculations.
	BollingerBands *BollingerBands[T]
}

// NewPercentB function initializes a new %B instance with the default parameters.
func NewPercentB[T helper.Number]() *PercentB[T] {
	return NewPercentBWithPeriod[T](DefaultBollingerBandsPeriod)
}

// NewPercentBWithPeriod function initializes a new %B instance with the given period.
func NewPercentBWithPeriod[T helper.Number](period int) *PercentB[T] {
	return &PercentB[T]{
		BollingerBands: NewBollingerBandsWithPeriod[T](period),
	}
}

// Compute function takes a channel of numbers and computes the %B over the specified period.
func (p *PercentB[T]) Compute(closings <-chan T) <-chan T {
	closingsSplice := helper.Duplicate(closings, 2)

	// Compute the Bollinger Bands
	upperBands, middleBands, lowerBands := p.BollingerBands.Compute(closingsSplice[0])

	// Skip closings until the Bollinger Bands are available
	closingsSplice[1] = helper.Skip(closingsSplice[1], p.BollingerBands.IdlePeriod())

	// Drain the middle bands
	go helper.Drain(middleBands)

	return helper.Operate3(upperBands, lowerBands, closingsSplice[1], func(upperBand, lowerBand, closing T) T {
		// %B = (Close - Lower Band) / (Upper Band - Lower Band)
		return (closing - lowerBand) / (upperBand - lowerBand)
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
