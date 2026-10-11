// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestCci(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		Cci   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/cci.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	low := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closing := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Cci })

	cci := trend.NewCci[float64]()

	actual := cci.Compute(high, low, closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, cci.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCciCancellation(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	highs := make(chan float64)
	lows := make(chan float64)
	closings := make(chan float64)

	cci := trend.NewCci[float64]()
	actual := cci.ComputeWithContext(ctx, highs, lows, closings)

	cancel()

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	if _, ok := <-actual; ok {
		t.Fatal("Cci channel should be closed after cancellation")
	}
}

func TestCciString(t *testing.T) {
	expected := "CCI(20)"
	actual := trend.NewCci[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
