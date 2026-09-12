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
	// DefaultSharpeRatioPeriodsPerYear is the default number of return periods in a year, matching
	// the approximate number of trading days used to annualize a Sharpe Ratio computed from daily
	// outcomes.
	DefaultSharpeRatioPeriodsPerYear = 252

	// sharpeRatioMinStdDev is the smallest standard deviation treated as nonzero. A near-zero
	// standard deviation is dominated by floating-point rounding noise rather than actual risk, and
	// dividing by it would blow up into an arbitrarily large, meaningless Sharpe Ratio.
	sharpeRatioMinStdDev = 1e-9
)

// SharpeRatioWithContext computes the annualized Sharpe Ratio for the given stream of cumulative
// outcome values, as produced by OutcomeWithContext, supporting context cancellation.
//
//	Sharpe = Mean(periodReturns) / StdDev(periodReturns) * Sqrt(periodsPerYear)
//
// The risk-free rate is assumed to be zero. The outcomes channel is assumed to hold one cumulative
// return value per trading period (for example, one per daily snapshot), which is exactly what
// OutcomeWithContext produces. Per-period returns are derived from the change in the underlying
// equity curve (1 + outcome) between consecutive outcomes.
//
// Fewer than two outcome values, or a return series with zero (or floating-point-noise-level)
// variance, such as a strategy that never trades, yields a Sharpe Ratio of zero rather than
// dividing by a near-zero standard deviation.
func SharpeRatioWithContext(ctx context.Context, outcomes <-chan float64, periodsPerYear int) float64 {
	equity := helper.IncrementByWithContext(ctx, outcomes, 1.0)
	returns := helper.ChanToSlice(helper.ChangeRatioWithContext(ctx, equity, 1))

	meanReturn := helper.Mean(returns)
	stdDevReturn := helper.StdDev(returns)

	if stdDevReturn < sharpeRatioMinStdDev {
		return 0
	}

	return (meanReturn / stdDevReturn) * math.Sqrt(float64(periodsPerYear))
}

// SharpeRatio wraps SharpeRatioWithContext for backwards compatibility.
//
// Deprecated: Use SharpeRatioWithContext instead.
func SharpeRatio(outcomes <-chan float64, periodsPerYear int) float64 {
	return SharpeRatioWithContext(context.Background(), outcomes, periodsPerYear)
}
