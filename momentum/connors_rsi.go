// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"
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

// ComputeWithContext function takes a channel of closings numbers and computes the Connors RSI.
func (c *ConnorsRsi[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	cs := helper.DuplicateWithContext(ctx, closings, 3)

	cs[0] = helper.BufferedWithContext(ctx, cs[0], c.PercentRankPeriod)
	cs[1] = helper.BufferedWithContext(ctx, cs[1], c.PercentRankPeriod)
	cs[2] = helper.BufferedWithContext(ctx, cs[2], c.PercentRankPeriod)

	// Component 1: RSI on closing prices
	rsis := c.Rsi.ComputeWithContext(ctx, cs[0])

	// Component 2: RSI on streak length
	streaks := c.Streak.ComputeWithContext(ctx, cs[1])
	streakRsis := c.StreakRsi.ComputeWithContext(ctx, streaks)

	// Component 3: PercentRank of ROC
	rocs := c.Roc.ComputeWithContext(ctx, cs[2])
	percentRanks := helper.PercentRankWithContext(ctx, rocs, c.PercentRankPeriod)

	// The three components run in parallel over the same closing-price stream, but each
	// reaches its first value at a different offset. Skip the faster branches so all
	// three land on the same time index before combining them.
	rsiIdle, streakRsiIdle, percentRankIdle := c.componentIdlePeriods()
	maxIdle := c.IdlePeriod()

	rsis = helper.SkipWithContext(ctx, rsis, maxIdle-rsiIdle)
	streakRsis = helper.SkipWithContext(ctx, streakRsis, maxIdle-streakRsiIdle)
	percentRanks = helper.SkipWithContext(ctx, percentRanks, maxIdle-percentRankIdle)

	// Combine: average of three components
	result := helper.MultiplyByWithContext(ctx, helper.AddWithContext(ctx, helper.AddWithContext(ctx, rsis, streakRsis),
		percentRanks,
	),
		T(1)/T(3),
	)

	return result
}

// componentIdlePeriods returns the idle period of each of the three parallel components,
// measured from the start of the closings stream.
func (c *ConnorsRsi[T]) componentIdlePeriods() (rsiIdle, streakRsiIdle, percentRankIdle int) {
	rsiIdle = c.Rsi.IdlePeriod()
	streakRsiIdle = c.Streak.IdlePeriod() + c.StreakRsi.IdlePeriod()
	// PercentRank ranks a value once it has period-1 predecessors, so it becomes
	// idle one input sooner than a naive period-based count would suggest.
	percentRankIdle = c.Roc.IdlePeriod() + c.PercentRankPeriod - 1

	return rsiIdle, streakRsiIdle, percentRankIdle
}

// IdlePeriod is the initial period that Connors RSI won't yield any results. The three
// components run in parallel, not sequentially, so this is the max of their individual
// idle periods, not the sum.
func (c *ConnorsRsi[T]) IdlePeriod() int {
	rsiIdle, streakRsiIdle, percentRankIdle := c.componentIdlePeriods()

	idle := rsiIdle
	if streakRsiIdle > idle {
		idle = streakRsiIdle
	}
	if percentRankIdle > idle {
		idle = percentRankIdle
	}

	return idle
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

// ComputeWithContext function takes a channel of closings numbers and computes the streak length.
// Positive values indicate consecutive up closes, negative values indicate consecutive down closes.
func (s *Streak[T]) ComputeWithContext(ctx context.Context, closings <-chan T) <-chan T {
	// Get the change
	changes := helper.ChangeWithContext(ctx, closings, 1)

	// Calculate streak based on direction
	result := helper.MapWithContext(ctx, changes, func(change T) T {
		if change > T(0) {
			return T(1)
		} else if change < T(0) {
			return T(-1)
		}
		return T(0)
	})

	// Now calculate cumulative streak
	cumulative := helper.MapWithPreviousWithContext(ctx, result, func(prev, curr T) T {
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

// String is the string representation of the Streak.
func (s *Streak[T]) String() string {
	return "STREAK"
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (c *ConnorsRsi[T]) Compute(closings <-chan T) <-chan T {
	return c.ComputeWithContext(context.Background(), closings)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *Streak[T]) Compute(closings <-chan T) <-chan T {
	return s.ComputeWithContext(context.Background(), closings)
}
