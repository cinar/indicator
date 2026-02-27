// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"
	"math"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultChopPeriod is the default period for the Choppiness Index (CHOP).
	DefaultChopPeriod = 14
)

// Chop represents the configuration parameters for calculating the Choppiness Index (CHOP).
// It is a technical analysis indicator that measures the market's trendiness or choppiness.
//
//	CHOP = 100 * LOG10( SUM(ATR(1), n) / (MAX(High, n) - MIN(Low, n)) ) / LOG10(n)
type Chop[T helper.Number] struct {
	// Period is the period for the CHOP.
	Period int
}

// NewChop function initializes a new CHOP instance with the default parameters.
func NewChop[T helper.Number]() *Chop[T] {
	return NewChopWithPeriod[T](DefaultChopPeriod)
}

// NewChopWithPeriod function initializes a new CHOP instance with the given period.
func NewChopWithPeriod[T helper.Number](period int) *Chop[T] {
	return &Chop[T]{
		Period: period,
	}
}

// Compute function takes channels of highs, lows, and closings, and computes the CHOP over the specified period.
func (c *Chop[T]) Compute(highs, lows, closings <-chan T) <-chan T {
	highs2 := helper.Duplicate(highs, 2)
	lows2 := helper.Duplicate(lows, 2)

	// TR calculation
	// Use previous closing by skipping highs and lows by one.
	trHighs := helper.Skip(highs2[0], 1)
	trLows := helper.Skip(lows2[0], 1)

	tr := helper.Operate3(trHighs, trLows, closings, func(high, low, closing T) T {
		return T(math.Max(float64(high-low), math.Max(float64(high-closing), float64(closing-low))))
	})

	sumTr := trend.NewMovingSumWithPeriod[T](c.Period).Compute(tr)

	// MAX(High, n) and MIN(Low, n)
	// They should be aligned with the TR bars (starting from bar 1).
	rangeHighs := helper.Skip(highs2[1], 1)
	rangeLows := helper.Skip(lows2[1], 1)

	maxHigh := trend.NewMovingMaxWithPeriod[T](c.Period).Compute(rangeHighs)
	minLow := trend.NewMovingMinWithPeriod[T](c.Period).Compute(rangeLows)

	log10n := math.Log10(float64(c.Period))

	chop := helper.Operate3(sumTr, maxHigh, minLow, func(sum, max, min T) T {
		diff := float64(max - min)
		if diff == 0 {
			return 0
		}

		val := 100 * math.Log10(float64(sum)/diff) / log10n
		return T(val)
	})

	return chop
}

// IdlePeriod is the initial period that CHOP won't yield any results.
func (c *Chop[T]) IdlePeriod() int {
	return c.Period
}

// String function returns a string representation of the CHOP.
func (c *Chop[T]) String() string {
	return fmt.Sprintf("CHOP(%d)", c.Period)
}
