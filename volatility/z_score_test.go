// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestZScore(t *testing.T) {
	type Data struct {
		Close    float64 `header:"Close"`
		Expected float64 `header:"Expected"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/z_score.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Expected })

	z := volatility.NewZScore[float64]()
	actual := z.Compute(closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.RoundDigits(expected, 2)
	expected = helper.Skip(expected, z.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestZScoreString(t *testing.T) {
	z := volatility.NewZScoreWithPeriod[float64](14)
	expected := "ZSCORE(14)"

	if z.String() != expected {
		t.Fatalf("expected %s actual %s", expected, z.String())
	}
}

func TestNewZScore(t *testing.T) {
	z := volatility.NewZScore[float64]()
	expectedPeriod := 20

	if z.Period != expectedPeriod {
		t.Fatalf("expected period %d actual %d", expectedPeriod, z.Period)
	}
}
