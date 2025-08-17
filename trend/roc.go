// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultRocPeriod is the default ROC period.
	DefaultRocPeriod = 9
)

// Roc represents the configuration parameters for calculating the Rate Of Change (ROC) indicator.
//
//	ROC = (Current Price - Price n periods ago) / Price n periods ago
type Roc[T helper.Float] struct {
	// Time period.
	Period int
}

// NewRoc function initializes a new Roc instance with the default parameters.
func NewRoc[T helper.Float]() *Roc[T] {
	return NewRocWithPeriod[T](DefaultRocPeriod)
}

// NewRocWithPeriod function initializes a new Roc instance with the given parameters.
func NewRocWithPeriod[T helper.Float](period int) *Roc[T] {
	if period <= 0 {
		period = DefaultRocPeriod
	}
	return &Roc[T]{
		Period: period,
	}
}

// Compute function takes a channel of numbers and computes the ROC and the signal line.
func (r *Roc[T]) Compute(values <-chan T) <-chan T {
	window := helper.NewRing[T](r.Period)

	rocs := helper.Map(values, func(value T) T {
		var result T

		if window.IsFull() {
			p, ok := window.Get()
			if ok && p != 0 {
				result = (value - p) / p
			}
		}
		window.Put(value)

		return result
	})

	rocs = helper.Skip(rocs, r.IdlePeriod())

	return rocs
}

// IdlePeriod is the initial period that ROC won't yield any results.
func (r *Roc[T]) IdlePeriod() int {
	return r.Period
}

// String is the string representation of the ROC.
func (r *Roc[T]) String() string {
	return fmt.Sprintf("ROC(%d)", r.Period)
}
