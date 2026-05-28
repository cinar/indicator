// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultHistoricalVolatilityPeriod is the default period for Historical Volatility (HV).
	DefaultHistoricalVolatilityPeriod = 21
)

// HistoricalVolatility represents the configuration parameters for calculating Historical Volatility (HV).
//
//	HV = StdDev(R_t, n)
//	where R_t = (P_t / P_(t-1)) - 1
//
// Refactored to utilize composition of helper.ChangeRatio and MovingStd.
type HistoricalVolatility[T helper.Number] struct {
	// Time period.
	Period int
}

// NewHistoricalVolatility function initializes a new Historical Volatility instance with the default parameters.
func NewHistoricalVolatility[T helper.Number]() *HistoricalVolatility[T] {
	return NewHistoricalVolatilityWithPeriod[T](DefaultHistoricalVolatilityPeriod)
}

// NewHistoricalVolatilityWithPeriod function initializes a new Historical Volatility instance with the given period.
func NewHistoricalVolatilityWithPeriod[T helper.Number](period int) *HistoricalVolatility[T] {
	if period <= 0 {
		period = DefaultHistoricalVolatilityPeriod
	}

	return &HistoricalVolatility[T]{
		Period: period,
	}
}

// ComputeWithContext function takes a channel of prices and computes the Historical Volatility over the specified period.
func (h *HistoricalVolatility[T]) ComputeWithContext(ctx context.Context, prices <-chan T) <-chan T {
	returns := helper.ChangeRatioWithContext(ctx, prices, 1)
	return NewMovingStdWithPeriod[T](h.Period).ComputeWithContext(ctx, returns)
}

// IdlePeriod is the initial period that Historical Volatility won't yield any results.
func (h *HistoricalVolatility[T]) IdlePeriod() int {
	// One bar for return series start and period-1 bars for window fill.
	return h.Period
}

// String function returns a string representation of the Historical Volatility.
func (h *HistoricalVolatility[T]) String() string {
	return fmt.Sprintf("HV(%d)", h.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (h *HistoricalVolatility[T]) Compute(prices <-chan T) <-chan T {
	return h.ComputeWithContext(context.Background(), prices)
}
