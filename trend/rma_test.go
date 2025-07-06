// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestRma(t *testing.T) {
	type Data struct {
		Close float64
		Rma   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/rma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Rma })

	rma := trend.NewRma[float64]()
	rma.Period = 15

	actual := rma.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, rma.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
