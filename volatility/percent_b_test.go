// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestPercentB(t *testing.T) {
	type Data struct {
		Close    float64
		PercentB float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/percent_b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.PercentB })

	percentB := volatility.NewPercentB[float64]()

	actual := percentB.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, percentB.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentBString(t *testing.T) {
	expected := "%B(10)"
	actual := volatility.NewPercentBWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
