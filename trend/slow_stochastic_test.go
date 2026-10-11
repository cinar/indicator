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

func TestSlowStochastic(t *testing.T) {
	type Data struct {
		Close float64 `header:"Close"`
		K     float64 `header:"K"`
		D     float64 `header:"D"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/slow_stochastic.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	s := trend.NewSlowStochastic[float64]()
	actualK, actualD := s.Compute(closing)

	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	inputs[1] = helper.Skip(inputs[1], s.IdlePeriod())

	for data := range inputs[1] {
		ak := <-actualK
		ad := <-actualD

		if ak != data.K {
			t.Fatalf("K: actual %v expected %v", ak, data.K)
		}

		if ad != data.D {
			t.Fatalf("D: actual %v expected %v", ad, data.D)
		}
	}
}

func TestSlowStochasticFlatMarket(t *testing.T) {
	s := trend.NewSlowStochastic[float64]()

	closing := make([]float64, s.IdlePeriod()+5)
	for i := range closing {
		closing[i] = 100
	}

	actualK, actualD := s.Compute(helper.SliceToChan(closing))

	// %K and %D share an upstream fan-out that requires both output
	// channels to be drained concurrently, so they are collected in
	// parallel rather than one at a time.
	var k, d []float64
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		k = helper.ChanToSlice(actualK)
	}()

	go func() {
		defer wg.Done()
		d = helper.ChanToSlice(actualD)
	}()

	wg.Wait()

	if len(k) == 0 {
		t.Fatal("expected at least one Slow Stochastic value")
	}

	for i := range k {
		if k[i] != 50 {
			t.Fatalf("expected Slow %%K of 50 for flat market, got %v", k[i])
		}

		if d[i] != 50 {
			t.Fatalf("expected Slow %%D of 50 for flat market, got %v", d[i])
		}
	}
}

func TestNewSlowStochasticWithPeriod(t *testing.T) {
	s := trend.NewSlowStochasticWithPeriod[float64](14, 3, 3)
	if s.Period != 14 {
		t.Fatalf("expected period 14, got %d", s.Period)
	}
}

func TestSlowStochasticString(t *testing.T) {
	expected := "SLOWSTOCH(14,3,3)"
	actual := trend.NewSlowStochasticWithPeriod[float64](14, 3, 3).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
