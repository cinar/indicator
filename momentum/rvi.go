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
	// DefaultRviPeriod is the default period for the Relative Vigor Index.
	DefaultRviPeriod = 10

	// DefaultRviSignalPeriod is the default signal line period for RVI.
	DefaultRviSignalPeriod = 4

	// RviFirPeriod is the FIR filter period (4 bars).
	RviFirPeriod = 4

	// RviFirSum is the sum of FIR weights (1+2+2+1 = 6).
	RviFirSum = 6
)

// Rvi represents the configuration parameters for calculating the
// Relative Vigor Index (RVI). The RVI is a momentum indicator that
// measures the strength of a trend by comparing close and open prices.
//
// The indicator uses a 4-bar FIR filter with weights 1-2-2-1:
//
//	Numerator = Close - Open
//	Denominator = High - Low
//	FIR(Numerator) = (1*Num[0] + 2*Num[1] + 2*Num[2] + 1*Num[3]) / 6
//	FIR(Denominator) = (1*Den[0] + 2*Den[1] + 2*Den[2] + 1*Den[3]) / 6
//	RVI = SMA(FIR(Numerator), period) / SMA(FIR(Denominator), period)
//	Signal = SMA(RVI, signalPeriod)
//
// Example:
//
//	rvi := momentum.NewRvi[float64]()
//	rviResult, signalResult := rvi.Compute(opens, highs, lows, closings)
type Rvi[T helper.Float] struct {
	// Period is the lookback period for the RVI.
	Period int

	// SignalPeriod is the signal line period.
	SignalPeriod int
}

// NewRvi function initializes a new RVI instance.
func NewRvi[T helper.Float]() *Rvi[T] {
	return &Rvi[T]{
		Period:       DefaultRviPeriod,
		SignalPeriod: DefaultRviSignalPeriod,
	}
}

// computeFir applies a 4-bar FIR filter with weights 1-2-2-1.
func computeFir[T helper.Float](c <-chan T) <-chan T {
	// FIR with weights 1-2-2-1:
	// output[n] = (1*input[n] + 2*input[n-1] + 2*input[n-2] + 1*input[n-3]) / 6

	// Duplicate to get delayed versions
	cs := helper.Duplicate(c, 4)

	// Shift each copy to get delayed values
	delayed0 := cs[0] // current
	delayed1 := helper.Shift(cs[1], 1, 0)
	delayed2 := helper.Shift(cs[2], 2, 0)
	delayed3 := helper.Shift(cs[3], 3, 0)

	// Apply weights: 1*current + 2*prev1 + 2*prev2 + 1*prev3
	weighted := helper.Add(
		helper.Add(delayed0, helper.MultiplyBy(delayed1, 2)),
		helper.Add(helper.MultiplyBy(delayed2, 2), delayed3),
	)

	// Divide by sum of weights (6)
	result := helper.MultiplyBy(weighted, T(1)/T(RviFirSum))

	// Skip first 3 values (FIR warmup)
	return helper.Skip(result, RviFirPeriod-1)
}

// Compute function takes channels of OHLC numbers and computes the
// Relative Vigor Index and its signal line.
func (r *Rvi[T]) Compute(opens, highs, lows, closings <-chan T) (rviResult <-chan T, signalResult <-chan T) {
	return r.computeSimple(opens, highs, lows, closings)
}

// computeSimple is a simpler implementation.
func (r *Rvi[T]) computeSimple(opens, highs, lows, closings <-chan T) (rviResult <-chan T, signalResult <-chan T) {
	// Collect inputs to allow multiple passes
	openVals := helper.ChanToSlice(opens)
	highVals := helper.ChanToSlice(highs)
	lowVals := helper.ChanToSlice(lows)
	closeVals := helper.ChanToSlice(closings)

	// Create channels for each calculation
	openChan := helper.SliceToChan(openVals)
	highChan := helper.SliceToChan(highVals)
	lowChan := helper.SliceToChan(lowVals)
	closeChan := helper.SliceToChan(closeVals)

	// Compute: Close - Open
	numeratorRaw := helper.Subtract(closeChan, openChan)

	// Compute: High - Low
	denominatorRaw := helper.Subtract(highChan, lowChan)

	// Apply 4-bar FIR filter
	numeratorFir := computeFir(numeratorRaw)
	denominatorFir := computeFir(denominatorRaw)

	// Apply SMA to filtered values
	smaNum := trend.NewSmaWithPeriod[T](r.Period)
	smaDen := trend.NewSmaWithPeriod[T](r.Period)

	smaNumerator := smaNum.Compute(numeratorFir)
	smaDenominator := smaDen.Compute(denominatorFir)

	// Divide: RVI = SMA(FIR(Numerator)) / SMA(FIR(Denominator))
	rviResult = helper.Divide(smaNumerator, smaDenominator)

	// Compute signal line
	signalSma := trend.NewSmaWithPeriod[T](r.SignalPeriod)
	signalResult = signalSma.Compute(rviResult)

	return rviResult, signalResult
}

// IdlePeriod is the initial period that RVI won't yield any results.
func (r *Rvi[T]) IdlePeriod() int {
	// FIR filter: RviFirPeriod-1 = 3
	// SMA: Period-1
	// Signal SMA: SignalPeriod-1
	// Total: 3 + (Period-1) + (SignalPeriod-1) = Period + SignalPeriod + 1
	return RviFirPeriod - 1 + r.Period - 1 + r.SignalPeriod - 1
}

// String is the string representation of the RVI.
func (r *Rvi[T]) String() string {
	return fmt.Sprintf("RVI(%d,%d)", r.Period, r.SignalPeriod)
}
