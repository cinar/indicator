// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"context"

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

// computeFir applies a 4-bar FIR filter with weights 1-2-2-1, propagating
// ctx through every stage so the filter closes its output promptly on
// cancellation instead of blocking on a stalled upstream channel.
func computeFir[T helper.Float](ctx context.Context, c <-chan T) <-chan T {
	// FIR with weights 1-2-2-1:
	// output[n] = (1*input[n] + 2*input[n-1] + 2*input[n-2] + 1*input[n-3]) / 6

	// Duplicate to get delayed versions
	cs := helper.DuplicateWithContext(ctx, c, 4)

	// Shift each copy to get delayed values
	delayed0 := cs[0] // current
	delayed1 := helper.ShiftWithContext(ctx, cs[1], 1, 0)
	delayed2 := helper.ShiftWithContext(ctx, cs[2], 2, 0)
	delayed3 := helper.ShiftWithContext(ctx, cs[3], 3, 0)

	// Apply weights: 1*current + 2*prev1 + 2*prev2 + 1*prev3
	weighted := helper.AddWithContext(ctx,
		helper.AddWithContext(ctx, delayed0, helper.MultiplyByWithContext(ctx, delayed1, 2)),
		helper.AddWithContext(ctx, helper.MultiplyByWithContext(ctx, delayed2, 2), delayed3),
	)

	// Divide by sum of weights (6)
	result := helper.MultiplyByWithContext(ctx, weighted, T(1)/T(RviFirSum))

	// Skip first 3 values (FIR warmup)
	return helper.SkipWithContext(ctx, result, RviFirPeriod-1)
}

// ComputeWithContext function takes channels of OHLC numbers and computes the
// Relative Vigor Index and its signal line.
func (r *Rvi[T]) ComputeWithContext(ctx context.Context, opens, highs, lows, closings <-chan T) (rviResult <-chan T, signalResult <-chan T) {
	return r.computeSimple(ctx, opens, highs, lows, closings)
}

// computeSimple is a simpler implementation.
//
// It used to collect the OHLC channels into slices with helper.ChanToSlice
// before streaming them back through helper.SliceToChan. That collection
// step ignored ctx and read its source channel to completion with a plain
// `for n := range c`, so on a channel that never closes (a genuinely
// open-ended stream) the call would block forever with no way to cancel
// it. Since each of opens/highs/lows/closings is only ever read once
// below (there was never an actual need for the "multiple passes" the
// collect step was collecting for), the fix removes the collect/replay
// round-trip entirely and wires the *WithContext helpers directly onto
// the input channels, matching how Fisher (momentum/fisher.go) streams.
// Every stage now selects on ctx.Done(), so ComputeWithContext returns
// immediately and both output channels close promptly as soon as ctx is
// cancelled, even if the inputs never close.
func (r *Rvi[T]) computeSimple(ctx context.Context, opens, highs, lows, closings <-chan T) (rviResult <-chan T, signalResult <-chan T) {
	// Compute: Close - Open
	numeratorRaw := helper.SubtractWithContext(ctx, closings, opens)

	// Compute: High - Low
	denominatorRaw := helper.SubtractWithContext(ctx, highs, lows)

	// Apply 4-bar FIR filter
	numeratorFir := computeFir(ctx, numeratorRaw)
	denominatorFir := computeFir(ctx, denominatorRaw)

	// Apply SMA to filtered values
	smaNum := trend.NewSmaWithPeriod[T](r.Period)
	smaDen := trend.NewSmaWithPeriod[T](r.Period)

	smaNumerator := smaNum.ComputeWithContext(ctx, numeratorFir)
	smaDenominator := smaDen.ComputeWithContext(ctx, denominatorFir)

	// Divide: RVI = SMA(FIR(Numerator)) / SMA(FIR(Denominator))
	rvi := helper.DivideWithContext(ctx, smaNumerator, smaDenominator)

	rviSplice := helper.DuplicateWithContext(ctx, rvi, 2)

	// Compute signal line
	signalSma := trend.NewSmaWithPeriod[T](r.SignalPeriod)
	signalResult = signalSma.ComputeWithContext(ctx, rviSplice[0])

	return helper.SkipWithContext(ctx, rviSplice[1], r.SignalPeriod-1), signalResult
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (r *Rvi[T]) Compute(opens, highs, lows, closings <-chan T) (rviResult <-chan T, signalResult <-chan T) {
	return r.ComputeWithContext(context.Background(), opens, highs, lows, closings)
}
