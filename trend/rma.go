// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import "github.com/cinar/indicator/v2/helper"

const (
	// DefaultRmaPeriod is the default RMA period.
	DefaultRmaPeriod = 20
)

// Rma represents the parameters for calculating Rolling Moving Average (RMA).
//
//	R[0] to R[p-1] is SMA(values)
//	R[p] and after is R[i] = ((R[i-1]*(p-1)) + v[i]) / p
//
// Example:
//
//	rma := trend.NewRma[float64]()
//	rma.Period = 10
//
//	result := rma.Compute(c)
type Rma[T helper.Number] struct {
	// Time period.
	Period int
}

// NewRma function initializes a new RMA instance with the default parameters.
func NewRma[T helper.Number]() *Rma[T] {
	return NewRmaWithPeriod[T](DefaultRmaPeriod)
}

// NewRmaWithPeriod function initializes a new RMA instance with the given period.
func NewRmaWithPeriod[T helper.Number](period int) *Rma[T] {
	return &Rma[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the RMA over the specified period.
func (r *Rma[T]) Compute(c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Initial RMA value is the SMA.
		sma := NewSma[T]()
		sma.Period = r.Period

		before := <-sma.Compute(helper.Head(c, r.Period))
		result <- before

		for n := range c {
			before = ((before * T(r.Period-1)) + n) / T(r.Period)
			result <- before
		}
	}()

	return result
}

// IdlePeriod is the initial period that RMA won't yield any results.
func (r *Rma[T]) IdlePeriod() int {
	return r.Period - 1
}
