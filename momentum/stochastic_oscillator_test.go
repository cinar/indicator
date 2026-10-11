// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestStochasticOscillator(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		K     float64
		D     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic_oscillator.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedK := helper.Map(inputs[3], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[4], func(d *Data) float64 { return d.D })

	so := momentum.NewStochasticOscillator[float64]()
	actualK, actualD := so.Compute(highs, lows, closings)
	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	expectedK = helper.Skip(expectedK, so.IdlePeriod())
	expectedD = helper.Skip(expectedD, so.IdlePeriod())

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD)
	if err != nil {
		t.Fatal(err)
	}
}

// TestStochasticOscillatorDifferentMaxMinPeriods verifies that StochasticOscillator correctly
// aligns the Max and Min branches when they are independently configured with different periods,
// rather than assuming Max.IdlePeriod() == Min.IdlePeriod(). Expected values were derived by hand
// (and cross-checked with an independent reference calculation) for this small synthetic input.
func TestStochasticOscillatorDifferentMaxMinPeriods(t *testing.T) {
	highs := helper.SliceToChan([]float64{10, 12, 11, 15, 14, 16})
	lows := helper.SliceToChan([]float64{5, 6, 4, 7, 8, 9})
	closings := helper.SliceToChan([]float64{8, 9, 7, 10, 11, 12})

	so := momentum.NewStochasticOscillator[float64]()
	so.Max.Period = 2
	so.Min.Period = 3

	if so.IdlePeriod() != 4 {
		t.Fatalf("actual idle period %v expected %v", so.IdlePeriod(), 4)
	}

	actualK, actualD := so.Compute(highs, lows, closings)
	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	expectedK := helper.SliceToChan([]float64{63.64, 55.56})
	expectedD := helper.SliceToChan([]float64{51.89, 57.91})

	err := helper.CheckEquals(actualK, expectedK, actualD, expectedD)
	if err != nil {
		t.Fatal(err)
	}
}

// TestStochasticOscillatorFlatMarket verifies that StochasticOscillator returns the neutral %K/%D
// of 50 instead of NaN when the high-low range is zero (a flat market with no price movement at
// all within the window).
func TestStochasticOscillatorFlatMarket(t *testing.T) {
	so := momentum.NewStochasticOscillator[float64]()

	const bars = 20
	flat := make([]float64, so.IdlePeriod()+bars)
	for i := range flat {
		flat[i] = 100
	}

	inputs := helper.Duplicate(helper.SliceToChan(flat), 3)

	actualK, actualD := so.Compute(inputs[0], inputs[1], inputs[2])

	expected := make([]float64, bars)
	for i := range expected {
		expected[i] = 50
	}

	err := helper.CheckEquals(actualK, helper.SliceToChan(expected), actualD, helper.SliceToChan(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStochasticOscillatorString(t *testing.T) {
	expected := "STOCHOSC(14,3)"
	actual := momentum.NewStochasticOscillator[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
