// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestTrueRange(t *testing.T) {
	type Data struct {
		High  float64 `header:"High"`
		Low   float64 `header:"Low"`
		Close float64 `header:"Close"`
		Tr    float64 `header:"Tr"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/tr.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Tr })

	tr := volatility.NewTrueRange[float64]()

	if tr.IdlePeriod() != 1 {
		t.Fatalf("expected IdlePeriod to be 1, got %d", tr.IdlePeriod())
	}

	if tr.String() != "TR" {
		t.Fatalf("expected String to be TR, got %s", tr.String())
	}

	actual := tr.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, tr.IdlePeriod())
	expected = helper.RoundDigits(expected, 2)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
