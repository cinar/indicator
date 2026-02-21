// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultKstRocPeriod1 is the default first ROC period.
	DefaultKstRocPeriod1 = 10

	// DefaultKstRocPeriod2 is the default second ROC period.
	DefaultKstRocPeriod2 = 15

	// DefaultKstRocPeriod3 is the default third ROC period.
	DefaultKstRocPeriod3 = 20

	// DefaultKstRocPeriod4 is the default fourth ROC period.
	DefaultKstRocPeriod4 = 30

	// DefaultKstSmaPeriod1 is the default first SMA period for smoothing.
	DefaultKstSmaPeriod1 = 10

	// DefaultKstSmaPeriod2 is the default second SMA period for smoothing.
	DefaultKstSmaPeriod2 = 10

	// DefaultKstSmaPeriod3 is the default third SMA period for smoothing.
	DefaultKstSmaPeriod3 = 10

	// DefaultKstSmaPeriod4 is the default fourth SMA period for smoothing.
	DefaultKstSmaPeriod4 = 15

	// DefaultKstSignalPeriod is the default signal line period.
	DefaultKstSignalPeriod = 9
)

// Kst represents the configuration parameters for calculating the
// Know Sure Thing (KST) oscillator. KST is a momentum oscillator
// based on the smoothed rate-of-change for four different timeframes.
// A KST value crossing above zero suggests a bullish trend, while
// crossing below zero indicates a bearish trend.
//
//	RCMA1 = SMA(ROC(close, rocPeriod1), smaPeriod1)
//	RCMA2 = SMA(ROC(close, rocPeriod2), smaPeriod2)
//	RCMA3 = SMA(ROC(close, rocPeriod3), smaPeriod3)
//	RCMA4 = SMA(ROC(close, rocPeriod4), smaPeriod4)
//	KST = (RCMA1 × 1) + (RCMA2 × 2) + (RCMA3 × 3) + (RCMA4 × 4)
//	Signal = SMA(KST, signalPeriod)
//
// Example:
//
//	kst := trend.NewKst[float64]()
//	kst.RocPeriod1 = 10
//	kst.RocPeriod2 = 15
//	kst.RocPeriod3 = 20
//	kst.RocPeriod4 = 30
//	kst.SmaPeriod1 = 10
//	kst.SmaPeriod2 = 10
//	kst.SmaPeriod3 = 10
//	kst.SmaPeriod4 = 15
//	kst.SignalPeriod = 9
//
//	kstResult, signalResult := kst.Compute(c)
type Kst[T helper.Float] struct {
	// RocPeriod1 is the first ROC period.
	RocPeriod1 int

	// RocPeriod2 is the second ROC period.
	RocPeriod2 int

	// RocPeriod3 is the third ROC period.
	RocPeriod3 int

	// RocPeriod4 is the fourth ROC period.
	RocPeriod4 int

	// SmaPeriod1 is the first SMA period for smoothing.
	SmaPeriod1 int

	// SmaPeriod2 is the second SMA period for smoothing.
	SmaPeriod2 int

	// SmaPeriod3 is the third SMA period for smoothing.
	SmaPeriod3 int

	// SmaPeriod4 is the fourth SMA period for smoothing.
	SmaPeriod4 int

	// SignalPeriod is the signal line period.
	SignalPeriod int
}

// NewKst function initializes a new KST instance with default parameters.
func NewKst[T helper.Float]() *Kst[T] {
	return &Kst[T]{
		RocPeriod1:   DefaultKstRocPeriod1,
		RocPeriod2:   DefaultKstRocPeriod2,
		RocPeriod3:   DefaultKstRocPeriod3,
		RocPeriod4:   DefaultKstRocPeriod4,
		SmaPeriod1:   DefaultKstSmaPeriod1,
		SmaPeriod2:   DefaultKstSmaPeriod2,
		SmaPeriod3:   DefaultKstSmaPeriod3,
		SmaPeriod4:   DefaultKstSmaPeriod4,
		SignalPeriod: DefaultKstSignalPeriod,
	}
}

// Compute function takes a channel of numbers and computes the KST
// and the signal line.
func (k *Kst[T]) Compute(c <-chan T) (kstResult <-chan T, signalResult <-chan T) {
	rocPeriods := []int{k.RocPeriod1, k.RocPeriod2, k.RocPeriod3, k.RocPeriod4}
	smaPeriods := []int{k.SmaPeriod1, k.SmaPeriod2, k.SmaPeriod3, k.SmaPeriod4}

	maxRocPeriod := 0
	for _, p := range rocPeriods {
		if p > maxRocPeriod {
			maxRocPeriod = p
		}
	}

	c = helper.Buffered(c, maxRocPeriod)

	cs := helper.Duplicate(c, 4)

	maxIdle := 0
	idles := make([]int, 4)
	for i := 0; i < 4; i++ {
		idles[i] = rocPeriods[i] + smaPeriods[i] - 1
		if idles[i] > maxIdle {
			maxIdle = idles[i]
		}
	}

	rcma := make([]<-chan T, 4)
	for i := 0; i < 4; i++ {
		roc := NewRocWithPeriod[T](rocPeriods[i])
		rcma[i] = roc.Compute(cs[i])

		sma := NewSmaWithPeriod[T](smaPeriods[i])
		rcma[i] = sma.Compute(rcma[i])

		skipCount := maxIdle - idles[i]
		if skipCount > 0 {
			rcma[i] = helper.Skip(rcma[i], skipCount)
		}
	}

	kst := helper.Add(
		helper.Add(
			helper.MultiplyBy(rcma[0], T(1)),
			helper.MultiplyBy(rcma[1], T(2)),
		),
		helper.Add(
			helper.MultiplyBy(rcma[2], T(3)),
			helper.MultiplyBy(rcma[3], T(4)),
		),
	)

	signal := NewSmaWithPeriod[T](k.SignalPeriod)
	signalResult = signal.Compute(kst)

	return kst, signalResult
}

// IdlePeriod is the initial period that KST won't yield any results.
func (k *Kst[T]) IdlePeriod() int {
	rocPeriods := []int{k.RocPeriod1, k.RocPeriod2, k.RocPeriod3, k.RocPeriod4}
	smaPeriods := []int{k.SmaPeriod1, k.SmaPeriod2, k.SmaPeriod3, k.SmaPeriod4}

	maxIdle := 0
	for i := 0; i < 4; i++ {
		idle := rocPeriods[i] + smaPeriods[i] - 1
		if idle > maxIdle {
			maxIdle = idle
		}
	}

	return maxIdle + k.SignalPeriod - 1
}

// String is the string representation of the KST.
func (k *Kst[T]) String() string {
	return fmt.Sprintf("KST(%d,%d,%d,%d,%d,%d,%d,%d,%d)",
		k.RocPeriod1, k.RocPeriod2, k.RocPeriod3, k.RocPeriod4,
		k.SmaPeriod1, k.SmaPeriod2, k.SmaPeriod3, k.SmaPeriod4,
		k.SignalPeriod)
}
