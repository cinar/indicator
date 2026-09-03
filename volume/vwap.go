// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultVwapPeriod is the default period for the VWAP.
	DefaultVwapPeriod = 14
)

// Vwap holds configuration parameters for calculating the Volume Weighted Average Price (VWAP). It provides the
// average price the asset has traded.
//
//	VWAP = Sum(Closing * Volume) / Sum(Volume)
//
// When no trading occurred anywhere in the window, the volume sum is 0 and a
// volume-weighted price is undefined (0/0). Unlike MFM/CMF, 0 is not a safe
// stand-in here: a VWAP of 0 would read as a real, implausibly low price rather
// than "no data," which is actively misleading if plotted or compared against
// actual prices (see the example WeightedAveragePriceStrategy, which crosses
// closing price against VWAP - a fabricated 0 would falsely signal a crossover
// every time). Instead, VWAP carries forward the last period with actual volume,
// which is the conventional real-market handling for an illiquid bar. Before any
// window has had volume, there is no prior value to carry forward, so VWAP
// returns the zero value of T.
//
// Example:
//
//	vwap := volume.NewVwap[float64]()
//	result := vwap.Compute(closings, volumes)
type Vwap[T helper.Float] struct {
	// Sum is the Moving Sum instance.
	Sum *trend.MovingSum[T]
}

// NewVwap function initializes a new VWAP instance with the default parameters.
func NewVwap[T helper.Float]() *Vwap[T] {
	return NewVwapWithPeriod[T](DefaultVwapPeriod)
}

// NewVwapWithPeriod function initializes a new VWAP instance with the given period.
func NewVwapWithPeriod[T helper.Float](period int) *Vwap[T] {
	return &Vwap[T]{
		Sum: trend.NewMovingSumWithPeriod[T](period),
	}
}

// ComputeWithContext function takes a channel of numbers and computes the VWAP.
func (v *Vwap[T]) ComputeWithContext(ctx context.Context, closings, volumes <-chan T) <-chan T {
	volumesSplice := helper.DuplicateWithContext(ctx, volumes, 2)

	sumPriceVolume := v.Sum.ComputeWithContext(ctx, helper.MultiplyWithContext(ctx, closings,
		volumesSplice[0],
	))
	sumVolume := v.Sum.ComputeWithContext(ctx, volumesSplice[1])

	// last holds the most recently computed valid VWAP, carried forward when a
	// window has no volume. OperateWithContext runs sequentially, so this
	// closure state is safe without synchronization.
	var last T

	return helper.OperateWithContext(ctx, sumPriceVolume, sumVolume, func(priceVolumeSum, volumeSum T) T {
		if volumeSum == 0 {
			return last
		}

		last = priceVolumeSum / volumeSum

		return last
	})
}

// IdlePeriod is the initial period that VWAP won't yield any results.
func (v *Vwap[T]) IdlePeriod() int {
	return v.Sum.IdlePeriod()
}

// String is the string representation of the VWAP.
func (v *Vwap[T]) String() string {
	return fmt.Sprintf("VWAP(%d)", v.Sum.Period)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (v *Vwap[T]) Compute(closings, volumes <-chan T) <-chan T {
	return v.ComputeWithContext(context.Background(), closings, volumes)
}
