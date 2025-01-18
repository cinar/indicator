// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestTrimaWithOddPeriod(t *testing.T) {
	type Data struct {
		Close float64
		Trima float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/trima_odd.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Trima })

	trima := trend.NewTrima[float64]()
	trima.Period = 15

	actual := trima.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, trima.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTrimaWithEvenPeriod(t *testing.T) {
	type Data struct {
		Close float64
		Trima float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/trima_even.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Trima })

	trima := trend.NewTrima[float64]()
	trima.Period = 20

	actual := trima.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, trima.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
