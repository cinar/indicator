// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestRviSimple(t *testing.T) {
	// Create sample OHLC data - steadily increasing prices
	count := 50
	opens := make([]float64, count)
	highs := make([]float64, count)
	lows := make([]float64, count)
	closes := make([]float64, count)

	for i := 0; i < count; i++ {
		opens[i] = float64(100 + i)
		highs[i] = float64(105 + i)
		lows[i] = float64(95 + i)
		closes[i] = float64(102 + i)
	}

	openChan := helper.SliceToChan(opens)
	highChan := helper.SliceToChan(highs)
	lowChan := helper.SliceToChan(lows)
	closeChan := helper.SliceToChan(closes)

	rvi := momentum.NewRvi[float64]()
	rviResult, signalResult := rvi.Compute(openChan, highChan, lowChan, closeChan)

	var rviSlice []float64
	var signalSlice []float64

	for v := range rviResult {
		rviSlice = append(rviSlice, v)
		signalSlice = append(signalSlice, <-signalResult)
	}

	if len(rviSlice) == 0 {
		t.Fatal("RVI produced no output")
	}

	// Check for NaN/Inf
	for i, v := range rviSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("RVI at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("RVI: %d values, Signal: %d values", len(rviSlice), len(signalSlice))
}

func TestRviString(t *testing.T) {
	rvi := momentum.NewRvi[float64]()
	expected := "RVI(10,4)"
	actual := rvi.String()
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestRviIdlePeriod(t *testing.T) {
	rvi := momentum.NewRvi[float64]()
	// FIR(3) + SMA(9) + SignalSMA(3) = 15
	expected := 15
	actual := rvi.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
