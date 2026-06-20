// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSlope(t *testing.T) {
	type Data struct {
		Close float64 `header:"Close"`
		Slope float64 `header:"Slope"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/slope.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Slope })

	slope := NewSlope[float64]()
	actual := slope.Compute(closings)
	actual = helper.RoundDigits(actual, 8)

	expected = helper.Skip(expected, slope.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSlopeFallbackPeriod(t *testing.T) {
	slope := NewSlopeWithPeriod[float64](-17)

	if slope.Period != DefaultSlopePeriod {
		t.Fatal("expected period to be fallback to default value")
	}
}

func TestSlopeToStringAndIdlePeriod(t *testing.T) {
	slope := NewSlopeWithPeriod[float64](0)
	if slope.IdlePeriod() != DefaultSlopePeriod {
		t.Fatalf("unexpected IdlePeriod: %d", slope.IdlePeriod())
	}

	slope.Period = 3
	if s := slope.String(); s != "SLOPE(3)" {
		t.Fatalf("unexpected String(): %s", s)
	}
}
