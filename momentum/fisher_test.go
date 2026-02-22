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

func TestFisherSimple(t *testing.T) {
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}

	input := helper.SliceToChan(prices)
	fisher := momentum.NewFisher[float64]()
	result := fisher.Compute(input)

	resultSlice := helper.ChanToSlice(result)

	if len(resultSlice) == 0 {
		t.Fatal("Fisher produced no output")
	}

	for i, v := range resultSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("Fisher at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("Fisher values: %v", resultSlice[:5])
}

func TestFisherString(t *testing.T) {
	fisher := momentum.NewFisher[float64]()
	expected := "Fisher(10)"
	actual := fisher.String()
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestFisherIdlePeriod(t *testing.T) {
	fisher := momentum.NewFisher[float64]()
	expected := 18
	actual := fisher.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestFisher(t *testing.T) {
	type Data struct {
		Close  float64
		Fisher float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/fisher.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	fisher := momentum.NewFisher[float64]()
	actual := fisher.Compute(closings)

	actual = helper.RoundDigits(actual, 2)

	inputs[1] = helper.Skip(inputs[1], fisher.IdlePeriod())

	err = helper.CheckEquals(
		helper.Map(actual, func(v float64) float64 { return v }),
		helper.Map(inputs[1], func(d *Data) float64 { return d.Fisher }),
	)
	if err != nil {
		t.Fatal(err)
	}
}
