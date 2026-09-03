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

func TestStochastic(t *testing.T) {
	type Data struct {
		Value float64
		K     float64
		D     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	values := helper.Map(inputs[0], func(d *Data) float64 { return d.Value })
	expectedK := helper.Map(inputs[1], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[2], func(d *Data) float64 { return d.D })

	s := trend.NewStochastic[float64]()
	actualK, actualD := s.Compute(values)
	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	expectedK = helper.Skip(expectedK, s.IdlePeriod())
	expectedD = helper.Skip(expectedD, s.IdlePeriod())

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStochasticFlatMarket(t *testing.T) {
	s := trend.NewStochastic[float64]()

	values := make([]float64, s.IdlePeriod()+5)
	for i := range values {
		values[i] = 100
	}

	actualK, actualD := s.Compute(helper.SliceToChan(values))

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
		t.Fatal("expected at least one Stochastic value")
	}

	for i := range k {
		if k[i] != 50 {
			t.Fatalf("expected %%K of 50 for flat market, got %v", k[i])
		}

		if d[i] != 50 {
			t.Fatalf("expected %%D of 50 for flat market, got %v", d[i])
		}
	}
}

func TestStochasticString(t *testing.T) {
	expected := "STOCHASTIC(10,3)"
	actual := trend.NewStochastic[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
