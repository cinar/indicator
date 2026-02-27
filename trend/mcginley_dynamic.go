// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"
	"math"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultMcGinleyDynamicPeriod is the default period for the McGinley Dynamic.
	DefaultMcGinleyDynamicPeriod = 14
)

// McGinleyDynamic represents the parameters for calculating the McGinley Dynamic.
// It is a technical analysis indicator that is an improvement over the Exponential
// Moving Average (EMA). It is designed to adjust for changes in market speed.
//
//	MD_today = MD_yesterday + (Close - MD_yesterday) / (Period * (Close / MD_yesterday)^4)
//
// Example:
//
//	md := trend.NewMcGinleyDynamic[float64]()
//	result := md.Compute(c)
type McGinleyDynamic[T helper.Number] struct {
	// Period is the smoothing period.
	Period int
}

// NewMcGinleyDynamic function initializes a new McGinley Dynamic instance with the default parameters.
func NewMcGinleyDynamic[T helper.Number]() *McGinleyDynamic[T] {
	return NewMcGinleyDynamicWithPeriod[T](DefaultMcGinleyDynamicPeriod)
}

// NewMcGinleyDynamicWithPeriod function initializes a new McGinley Dynamic instance with the given period.
func NewMcGinleyDynamicWithPeriod[T helper.Number](period int) *McGinleyDynamic[T] {
	return &McGinleyDynamic[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the McGinley Dynamic over the specified period.
func (m *McGinleyDynamic[T]) Compute(c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		var before float64
		first := true

		for n := range c {
			val := float64(n)
			if first {
				before = val
				first = false
				result <- T(before)
				continue
			}

			if before == 0 {
				before = val
				result <- T(before)
				continue
			}

			// MD_today = MD_yesterday + (Close - MD_yesterday) / (Period * (Close / MD_yesterday)^4)
			ratio := val / before
			before = before + (val-before)/(float64(m.Period)*math.Pow(ratio, 4))
			result <- T(before)
		}
	}()

	return result
}

// IdlePeriod is the initial period that McGinley Dynamic yield any results.
func (m *McGinleyDynamic[T]) IdlePeriod() int {
	return 0
}

// String is the string representation of the McGinley Dynamic.
func (m *McGinleyDynamic[T]) String() string {
	return fmt.Sprintf("MD(%d)", m.Period)
}
