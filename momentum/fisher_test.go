// Copyright (c) 2021-2026 The Indicator Authors.
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
	prices := make([]float64, 400)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}

	input := helper.SliceToChan(prices)
	fisher := momentum.NewFisher[float64]()
	result := fisher.Compute(input)

	resultSlice := helper.ChanToSlice(result)

	expected := len(prices) - fisher.IdlePeriod()
	if len(resultSlice) != expected {
		t.Fatalf("expected %d values, got %d", expected, len(resultSlice))
	}

	for i, v := range resultSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("Fisher at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("Fisher values: %v", resultSlice[:5])
}

// referenceFisher computes the Fisher Transform with a plain sliding-window
// min/max over a slice, independently of Fisher's channel implementation.
func referenceFisher(closings []float64, period int) []float64 {
	var out []float64

	for i := period - 1; i < len(closings); i++ {
		window := closings[i-period+1 : i+1]

		lo, hi := window[0], window[0]
		for _, v := range window {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}

		x := 2*((closings[i]-lo)/(hi-lo)) - 1
		if x > momentum.FisherClamp {
			x = momentum.FisherClamp
		}
		if x < -momentum.FisherClamp {
			x = -momentum.FisherClamp
		}

		out = append(out, 0.5*math.Log((1+x)/(1-x)))
	}

	return out
}

func TestFisherValues(t *testing.T) {
	closings := make([]float64, 400)
	for i := range closings {
		closings[i] = 100 + 10*math.Sin(float64(i)/3)
	}

	fisher := momentum.NewFisher[float64]()
	out := helper.ChanToSlice(fisher.Compute(helper.SliceToChan(closings)))

	expected := referenceFisher(closings, fisher.Period)
	if len(out) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(out))
	}

	for i := range expected {
		if math.Abs(out[i]-expected[i]) > 1e-9 {
			t.Fatalf("value %d: expected %v, got %v", i, expected[i], out[i])
		}
	}
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
	expected := 9
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
