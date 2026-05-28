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
	// DefaultCoppockCurveRocPeriod1 is the default first ROC period.
	DefaultCoppockCurveRocPeriod1 = 14

	// DefaultCoppockCurveRocPeriod2 is the default second ROC period.
	DefaultCoppockCurveRocPeriod2 = 11

	// DefaultCoppockCurveWmaPeriod is the default WMA period.
	DefaultCoppockCurveWmaPeriod = 10
)

// CoppockCurve represents the configuration parameters for calculating the
// Coppock Curve oscillator. The Coppock Curve is a momentum indicator
// used to identify long-term buying opportunities in equity indices.
//
//	Coppock Curve = WMA(ROC(14) + ROC(11), 10)
//
// Example:
//
//	cc := momentum.NewCoppockCurve[float64]()
//	result := cc.Compute(closings)
type CoppockCurve[T helper.Float] struct {
	// RocPeriod1 is the first ROC period.
	RocPeriod1 int

	// RocPeriod2 is the second ROC period.
	RocPeriod2 int

	// WmaPeriod is the WMA period for smoothing.
	WmaPeriod int
}

// NewCoppockCurve function initializes a new CoppockCurve instance with default parameters.
func NewCoppockCurve[T helper.Float]() *CoppockCurve[T] {
	return &CoppockCurve[T]{
		RocPeriod1: DefaultCoppockCurveRocPeriod1,
		RocPeriod2: DefaultCoppockCurveRocPeriod2,
		WmaPeriod:  DefaultCoppockCurveWmaPeriod,
	}
}

// NewCoppockCurveWithPeriods function initializes a new CoppockCurve instance with the given periods.
func NewCoppockCurveWithPeriods[T helper.Float](rocPeriod1, rocPeriod2, wmaPeriod int) *CoppockCurve[T] {
	if rocPeriod1 <= 0 {
		rocPeriod1 = DefaultCoppockCurveRocPeriod1
	}
	if rocPeriod2 <= 0 {
		rocPeriod2 = DefaultCoppockCurveRocPeriod2
	}
	if wmaPeriod <= 0 {
		wmaPeriod = DefaultCoppockCurveWmaPeriod
	}
	return &CoppockCurve[T]{
		RocPeriod1: rocPeriod1,
		RocPeriod2: rocPeriod2,
		WmaPeriod:  wmaPeriod,
	}
}

// Compute function takes a channel of closings and computes the Coppock Curve.
func (c *CoppockCurve[T]) Compute(values <-chan T) <-chan T {
	maxRocPeriod := c.RocPeriod1
	if c.RocPeriod2 > maxRocPeriod {
		maxRocPeriod = c.RocPeriod2
	}

	values = helper.Buffered(values, maxRocPeriod)
	valuesSplice := helper.Duplicate(values, 2)

	roc1 := trend.NewRocWithPeriod[T](c.RocPeriod1)
	roc1Values := helper.MultiplyBy(roc1.Compute(valuesSplice[0]), 100)

	roc2 := trend.NewRocWithPeriod[T](c.RocPeriod2)
	roc2Values := helper.MultiplyBy(roc2.Compute(valuesSplice[1]), 100)

	// Align ROC streams. Both Compute calls return streams that have already skipped their own IdlePeriod.
	// To align them to maxRocPeriod, we skip the remaining difference.
	if c.RocPeriod1 < maxRocPeriod {
		roc1Values = helper.Skip(roc1Values, maxRocPeriod-c.RocPeriod1)
	}
	if c.RocPeriod2 < maxRocPeriod {
		roc2Values = helper.Skip(roc2Values, maxRocPeriod-c.RocPeriod2)
	}

	sumRocs := helper.Add(roc1Values, roc2Values)

	wma := trend.NewWmaWith[T](c.WmaPeriod)
	return wma.Compute(sumRocs)
}

// IdlePeriod is the initial period that Coppock Curve won't yield any results.
func (c *CoppockCurve[T]) IdlePeriod() int {
	maxRocPeriod := c.RocPeriod1
	if c.RocPeriod2 > maxRocPeriod {
		maxRocPeriod = c.RocPeriod2
	}

	// maxRocPeriod (from ROC) + (WmaPeriod - 1) (from WMA)
	return maxRocPeriod + c.WmaPeriod - 1
}

// String is the string representation of the Coppock Curve.
func (c *CoppockCurve[T]) String() string {
	return fmt.Sprintf("CoppockCurve(%d,%d,%d)", c.RocPeriod1, c.RocPeriod2, c.WmaPeriod)
}
