// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultConnorsRsiRsiPeriod is the default RSI period.
	DefaultConnorsRsiRsiPeriod = 3
	// DefaultConnorsRsiStreakRsiPeriod is the default Streak RSI period.
	DefaultConnorsRsiStreakRsiPeriod = 2
	// DefaultConnorsRsiPercentRankPeriod is the default PercentRank period.
	DefaultConnorsRsiPercentRankPeriod = 100
)

// ConnorsRsi represents the configuration parameters for calculating the Connors RSI.
// It is a momentum indicator that combines three components:
// 1. RSI of closing prices
// 2. RSI of up/down streak length
// 3. Percentile rank of the rate of change
//
//	CRSI = (RSI(3) + RSI(Streak, 2) + PercentRank(ROC, 100)) / 3
//
// Example:
//
//	connorsRsi := momentum.NewConnorsRsi[float64]()
//	result := connorsRsi.Compute(closings)
type ConnorsRsi[T helper.Float] struct {
	// RsiPeriod is the period for the RSI on closing prices.
	RsiPeriod int
	// StreakRsiPeriod is the period for the RSI on streak length.
	StreakRsiPeriod int
	// PercentRankPeriod is the period for the PercentRank of ROC.
	PercentRankPeriod int

	// Rsi is the RSI instance for closing prices.
	Rsi *Rsi[T]
	// StreakRsi is the RSI instance for streak length.
	StreakRsi *Rsi[T]
	// Roc is the Rate of Change instance.
	Roc *trend.Roc[T]
	// Streak is the streak calculator instance.
	Streak *Streak[T]
}

// NewConnorsRsi function initializes a new Connors RSI instance with the default parameters.
func NewConnorsRsi[T helper.Float]() *ConnorsRsi[T] {
	return NewConnorsRsiWithPeriods[T](
		DefaultConnorsRsiRsiPeriod,
		DefaultConnorsRsiStreakRsiPeriod,
		DefaultConnorsRsiPercentRankPeriod,
	)
}

// NewConnorsRsiWithPeriods function initializes a new Connors RSI instance with the given periods.
func NewConnorsRsiWithPeriods[T helper.Float](rsiPeriod, streakRsiPeriod, percentRankPeriod int) *ConnorsRsi[T] {
	if rsiPeriod <= 0 {
		rsiPeriod = DefaultConnorsRsiRsiPeriod
	}
	if streakRsiPeriod <= 0 {
		streakRsiPeriod = DefaultConnorsRsiStreakRsiPeriod
	}
	if percentRankPeriod <= 0 {
		percentRankPeriod = DefaultConnorsRsiPercentRankPeriod
	}

	return &ConnorsRsi[T]{
		RsiPeriod:         rsiPeriod,
		StreakRsiPeriod:   streakRsiPeriod,
		PercentRankPeriod: percentRankPeriod,
		Rsi:               NewRsiWithPeriod[T](rsiPeriod),
		StreakRsi:         NewRsiWithPeriod[T](streakRsiPeriod),
		Roc:               trend.NewRocWithPeriod[T](1),
		Streak:            NewStreak[T](),
	}
}

// Compute function takes a channel of closings numbers and computes the Connors RSI.
func (c *ConnorsRsi[T]) Compute(closings <-chan T) <-chan T {
	cs := helper.Duplicate(closings, 3)

	cs[0] = helper.Buffered(cs[0], 100)
	cs[1] = helper.Buffered(cs[1], 100)
	cs[2] = helper.Buffered(cs[2], 100)

	// Component 1: RSI on closing prices
	rsis := c.Rsi.Compute(cs[0])

	// Component 2: RSI on streak length
	streaks := c.Streak.Compute(cs[1])
	streakRsis := c.StreakRsi.Compute(streaks)

	// Component 3: PercentRank of ROC
	rocs := c.Roc.Compute(cs[2])
	percentRanks := helper.PercentRank(rocs, c.PercentRankPeriod)

	// Combine: average of three components
	result := helper.MultiplyBy(
		helper.Add(
			helper.Add(rsis, streakRsis),
			percentRanks,
		),
		T(1)/T(3),
	)

	return result
}

// IdlePeriod is the initial period that Connors RSI won't yield any results.
func (c *ConnorsRsi[T]) IdlePeriod() int {
	// ROC period 1 + RSI period 3 + RMA period 14 + PercentRank period 100
	// = 1 + 3 + 14 + 100 = 118
	return c.Roc.IdlePeriod() + c.Rsi.IdlePeriod() + c.PercentRankPeriod
}

// String is the string representation of the Connors RSI.
func (c *ConnorsRsi[T]) String() string {
	return fmt.Sprintf("ConnorsRSI(%d, %d, %d)", c.RsiPeriod, c.StreakRsiPeriod, c.PercentRankPeriod)
}

// Streak represents the configuration for calculating the up/down streak length.
// The streak is the number of consecutive days the price has closed up or down.
type Streak[T helper.Float] struct{}

// NewStreak function initializes a new Streak instance.
func NewStreak[T helper.Float]() *Streak[T] {
	return &Streak[T]{}
}

// Compute function takes a channel of closings numbers and computes the streak length.
// Positive values indicate consecutive up closes, negative values indicate consecutive down closes.
func (s *Streak[T]) Compute(closings <-chan T) <-chan T {
	// Get the change
	changes := helper.Change(closings, 1)

	// Calculate streak based on direction
	result := helper.Map(changes, func(change T) T {
		if change > T(0) {
			return T(1)
		} else if change < T(0) {
			return T(-1)
		}
		return T(0)
	})

	// Now calculate cumulative streak
	cumulative := helper.MapWithPrevious(result, func(prev, curr T) T {
		if curr > T(0) {
			// Price went up - increment if previous was positive, else start at 1
			if prev > T(0) {
				return prev + T(1)
			}
			return T(1)
		} else if curr < T(0) {
			// Price went down - decrement if previous was negative, else start at -1
			if prev < T(0) {
				return prev - T(1)
			}
			return T(-1)
		}
		// Price unchanged - reset to 0
		return T(0)
	}, T(0))

	return cumulative
}

// IdlePeriod is the initial period that Streak won't yield any results.
func (s *Streak[T]) IdlePeriod() int {
	return 1
}
