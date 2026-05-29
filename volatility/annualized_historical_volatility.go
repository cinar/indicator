// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"
	"math"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultAnnualizedHistoricalVolatilityPeriod is the default period for
	// Annualized Historical Volatility (AHV).
	DefaultAnnualizedHistoricalVolatilityPeriod = 21

	// DefaultTradingDaysPerYear is the standard number of trading days in a year.
	DefaultTradingDaysPerYear = 252
)

// AnnualizedHistoricalVolatility represents the configuration parameters for calculating
// Annualized Historical Volatility (AHV). It annualizes the Historical Volatility by
// multiplying it by the square root of the number of trading days per year.
//
//	AHV = HV × √TradingDaysPerYear
type AnnualizedHistoricalVolatility[T helper.Number] struct {
	// Hv is the underlying Historical Volatility indicator.
	Hv *HistoricalVolatility[T]

	// TradingDaysPerYear is the number of trading days in a year (default: 252).
	TradingDaysPerYear int
}

// NewAnnualizedHistoricalVolatility function initializes a new Annualized Historical Volatility
// instance with the default parameters.
func NewAnnualizedHistoricalVolatility[T helper.Number]() *AnnualizedHistoricalVolatility[T] {
	return NewAnnualizedHistoricalVolatilityWithPeriod[T](DefaultAnnualizedHistoricalVolatilityPeriod)
}

// NewAnnualizedHistoricalVolatilityWithPeriod function initializes a new Annualized Historical
// Volatility instance with the given period.
func NewAnnualizedHistoricalVolatilityWithPeriod[T helper.Number](period int) *AnnualizedHistoricalVolatility[T] {
	return &AnnualizedHistoricalVolatility[T]{
		Hv:                 NewHistoricalVolatilityWithPeriod[T](period),
		TradingDaysPerYear: DefaultTradingDaysPerYear,
	}
}

// ComputeWithContext function takes a channel of prices and computes the Annualized Historical
// Volatility over the specified period.
func (a *AnnualizedHistoricalVolatility[T]) ComputeWithContext(ctx context.Context, prices <-chan T) <-chan T {
	hv := a.Hv.ComputeWithContext(ctx, prices)
	return helper.MultiplyByWithContext(ctx, hv, T(math.Sqrt(float64(a.TradingDaysPerYear))))
}

// IdlePeriod is the initial period that Annualized Historical Volatility won't yield any results.
func (a *AnnualizedHistoricalVolatility[T]) IdlePeriod() int {
	return a.Hv.IdlePeriod()
}

// String function returns a string representation of the Annualized Historical Volatility.
func (a *AnnualizedHistoricalVolatility[T]) String() string {
	return fmt.Sprintf("AHV(%d)", a.Hv.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *AnnualizedHistoricalVolatility[T]) Compute(prices <-chan T) <-chan T {
	return a.ComputeWithContext(context.Background(), prices)
}
