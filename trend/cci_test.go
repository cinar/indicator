// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestCci(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		Cci   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/cci.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	low := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closing := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Cci })

	cci := trend.NewCci[float64]()

	actual := cci.Compute(high, low, closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, cci.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
