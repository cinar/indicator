// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestUlcerIndex(t *testing.T) {
	type Data struct {
		Close      float64
		UlcerIndex float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ulcer_index.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.UlcerIndex })

	ui := volatility.NewUlcerIndex[float64]()
	actual := ui.Compute(closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, ui.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUlcerIndexZeroOrNegativePrice(t *testing.T) {
	ui := volatility.NewUlcerIndex[float64]()
	length := ui.IdlePeriod() + 20

	// A degenerate closing series that is zero or negative throughout, so
	// the rolling High Closings (a moving max) is never positive.
	closings := make([]float64, length)
	for i := range closings {
		if i%2 == 0 {
			closings[i] = 0
		} else {
			closings[i] = -1
		}
	}

	actual := helper.ChanToSlice(ui.Compute(helper.SliceToChan(closings)))

	if len(actual) == 0 {
		t.Fatal("expected at least one Ulcer Index value")
	}

	for i, v := range actual {
		if v != 0 {
			t.Fatalf("expected Ulcer Index of 0 for non-positive prices at %d, got %v", i, v)
		}
	}
}

func TestUlcerIndexString(t *testing.T) {
	expected := "ULCERINDEX(14)"
	actual := volatility.NewUlcerIndex[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
