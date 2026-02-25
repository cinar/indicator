// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestKvo(t *testing.T) {
	count := 100
	highs := make([]float64, count)
	lows := make([]float64, count)
	volumes := make([]float64, count)

	for i := 0; i < count; i++ {
		highs[i] = float64(100 + i)
		lows[i] = float64(90 + i)
		volumes[i] = float64(1000000 + i*1000)
	}

	highChan := helper.SliceToChan(highs)
	lowChan := helper.SliceToChan(lows)
	volChan := helper.SliceToChan(volumes)

	kvo := volume.NewKvo[float64]()
	kvoResult, signalResult := kvo.Compute(highChan, lowChan, volChan)

	var kvoSlice []float64
	var signalSlice []float64

	for v := range kvoResult {
		kvoSlice = append(kvoSlice, v)
		signalSlice = append(signalSlice, <-signalResult)
	}

	if len(kvoSlice) == 0 {
		t.Fatal("KVO produced no output")
	}

	for i, v := range kvoSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("KVO at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("KVO: %d values, Signal: %d values", len(kvoSlice), len(signalSlice))
}

func TestKvoIdlePeriod(t *testing.T) {
	kvo := volume.NewKvo[float64]()
	expected := 67
	actual := kvo.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
