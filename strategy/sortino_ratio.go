// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"context"
	"math"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultSortinoRatioPeriodsPerYear is the default number of return periods in a year, matching
	// the approximate number of trading days used to annualize a Sortino Ratio computed from daily
	// outcomes.
	DefaultSortinoRatioPeriodsPerYear = 252

	// sortinoRatioMinDownsideDeviation is the smallest downside deviation treated as nonzero. A
	// near-zero downside deviation is dominated by floating-point rounding noise rather than actual
	// downside risk, and dividing by it would blow up into an arbitrarily large, meaningless Sortino
	// Ratio.
	sortinoRatioMinDownsideDeviation = 1e-9

	// sortinoRatioMinimumAcceptableReturn is the per-period return below which a return counts as
	// downside risk.
	sortinoRatioMinimumAcceptableReturn = 0
)

// SortinoRatioWithContext computes the annualized Sortino Ratio for the given stream of cumulative
// outcome values, as produced by OutcomeWithContext, supporting context cancellation.
//
//	Sortino = Mean(periodReturns) / DownsideDeviation(periodReturns) * Sqrt(periodsPerYear)
//
// Unlike the Sharpe Ratio, which divides by the standard deviation of all returns, the Sortino
// Ratio divides by the downside deviation: the root-mean-square of only the shortfall below a
// minimum acceptable return (assumed to be zero here), with periods at or above that return
// contributing zero. Two return series with identical upside volatility but different downside
// volatility therefore yield different Sortino Ratios, even when their Sharpe Ratios are equal.
//
// The outcomes channel is assumed to hold one cumulative return value per trading period (for
// example, one per daily snapshot), which is exactly what OutcomeWithContext produces. Per-period
// returns are derived from the change in the underlying equity curve (1 + outcome) between
// consecutive outcomes.
//
// Fewer than two outcome values, or a return series with zero (or floating-point-noise-level)
// downside deviation, such as a strategy whose returns never fall below the minimum acceptable
// return, yields a Sortino Ratio of zero rather than dividing by a near-zero downside deviation.
func SortinoRatioWithContext(ctx context.Context, outcomes <-chan float64, periodsPerYear int) float64 {
	equity := helper.IncrementByWithContext(ctx, outcomes, 1.0)
	returns := helper.ChanToSlice(helper.ChangeRatioWithContext(ctx, equity, 1))

	meanReturn := helper.Mean(returns)
	downsideDeviationReturn := downsideDeviation(returns, sortinoRatioMinimumAcceptableReturn)

	if downsideDeviationReturn < sortinoRatioMinDownsideDeviation {
		return 0
	}

	return (meanReturn / downsideDeviationReturn) * math.Sqrt(float64(periodsPerYear))
}

// SortinoRatio wraps SortinoRatioWithContext for backwards compatibility.
//
// Deprecated: Use SortinoRatioWithContext instead.
func SortinoRatio(outcomes <-chan float64, periodsPerYear int) float64 {
	return SortinoRatioWithContext(context.Background(), outcomes, periodsPerYear)
}

// downsideDeviation returns the population downside deviation of the given returns relative to
// the given minimum acceptable return (MAR): the root-mean-square of the shortfall below MAR,
// treating a return at or above MAR as zero shortfall. Returns zero for an empty slice.
func downsideDeviation(returns []float64, mar float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	var sumSquaredShortfall float64

	for _, r := range returns {
		if shortfall := mar - r; shortfall > 0 {
			sumSquaredShortfall += shortfall * shortfall
		}
	}

	return math.Sqrt(sumSquaredShortfall / float64(len(returns)))
}
