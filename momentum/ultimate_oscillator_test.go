// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestUltimateOscillator(t *testing.T) {
	type Data struct {
		High  float64 `header:"High"`
		Low   float64 `header:"Low"`
		Close float64 `header:"Close"`
		Uo    float64 `header:"UO"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ultimate_oscillator.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Uo })

	uo := momentum.NewUltimateOscillator[float64]()
	actual := uo.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, uo.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// TestUltimateOscillatorFlatMarket verifies that UltimateOscillator returns its neutral value of
// 50 instead of NaN when there is no true range at all across the lookback windows (a flat
// market with no price movement).
func TestUltimateOscillatorFlatMarket(t *testing.T) {
	uo := momentum.NewUltimateOscillator[float64]()

	bars := uo.IdlePeriod() + 20
	flat := make([]float64, bars)
	for i := range flat {
		flat[i] = 100
	}

	inputs := helper.Duplicate(helper.SliceToChan(flat), 3)

	actual := helper.ChanToSlice(uo.Compute(inputs[0], inputs[1], inputs[2]))

	if len(actual) == 0 {
		t.Fatal("expected at least one UO value")
	}

	for _, v := range actual {
		if v != 50 {
			t.Fatalf("expected UO of 50 for flat market, got %v", v)
		}
	}
}
