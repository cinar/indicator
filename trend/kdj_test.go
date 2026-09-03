// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"sync"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestKdj(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		K     float64
		D     float64
		J     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/kdj.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 6)
	high := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	low := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closing := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedK := helper.Map(inputs[3], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[4], func(d *Data) float64 { return d.D })
	expectedJ := helper.Map(inputs[5], func(d *Data) float64 { return d.J })

	kdj := trend.NewKdj[float64]()
	actualK, actualD, actualJ := kdj.Compute(high, low, closing)

	actualK = helper.RoundDigits(actualK, 2)
	actualK = helper.Shift(actualK, kdj.IdlePeriod(), 0)

	actualD = helper.RoundDigits(actualD, 2)
	actualD = helper.Shift(actualD, kdj.IdlePeriod(), 0)

	actualJ = helper.RoundDigits(actualJ, 2)
	actualJ = helper.Shift(actualJ, kdj.IdlePeriod(), 0)

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD, actualJ, expectedJ)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKdjFlatMarket(t *testing.T) {
	count := trend.DefaultKdjMinMaxPeriod + trend.DefaultKdjSma1Period + trend.DefaultKdjSma2Period + 5
	high := make([]float64, count)
	low := make([]float64, count)
	closing := make([]float64, count)

	for i := range count {
		high[i] = 100
		low[i] = 100
		closing[i] = 100
	}

	kdj := trend.NewKdj[float64]()
	actualK, actualD, actualJ := kdj.Compute(helper.SliceToChan(high), helper.SliceToChan(low), helper.SliceToChan(closing))

	// K, D, and J share an upstream fan-out that requires all three output
	// channels to be drained concurrently, so they are collected in
	// parallel rather than one at a time.
	var k, d, j []float64
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		k = helper.ChanToSlice(actualK)
	}()

	go func() {
		defer wg.Done()
		d = helper.ChanToSlice(actualD)
	}()

	go func() {
		defer wg.Done()
		j = helper.ChanToSlice(actualJ)
	}()

	wg.Wait()

	if len(k) == 0 {
		t.Fatal("expected at least one KDJ value")
	}

	for i := range k {
		if k[i] != 50 {
			t.Fatalf("expected K of 50 for a zero-range market, got %v", k[i])
		}

		if d[i] != 50 {
			t.Fatalf("expected D of 50 for a zero-range market, got %v", d[i])
		}

		if j[i] != 50 {
			t.Fatalf("expected J of 50 for a zero-range market, got %v", j[i])
		}
	}
}

func TestKdjString(t *testing.T) {
	expected := "KDJ(9,3,3)"
	actual := trend.NewKdj[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
