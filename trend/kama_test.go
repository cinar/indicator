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

func TestKama(t *testing.T) {
	type Data struct {
		Close float64
		Kama  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/kama.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Kama })

	kama := trend.NewKama[float64]()
	actual := kama.Compute(closing)

	actual = helper.RoundDigits(actual, 2)
	expected = helper.Skip(expected, kama.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKamaEmpty(t *testing.T) {
	input := helper.SliceToChan([]float64{})
	expected := helper.SliceToChan([]float64{})

	kama := trend.NewKama[float64]()
	actual := kama.Compute(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKamaCancellation(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	closings := make(chan float64)

	kama := trend.NewKama[float64]()
	actual := kama.ComputeWithContext(ctx, closings)

	cancel()

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	if _, ok := <-actual; ok {
		t.Fatal("Kama channel should be closed after cancellation")
	}
}

func TestKamaString(t *testing.T) {
	expected := "KAMA(1,2,3)"

	kama := trend.NewKama[float64]()
	kama.ErPeriod = 1
	kama.FastScPeriod = 2
	kama.SlowScPeriod = 3

	actual := kama.String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
