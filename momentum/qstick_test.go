// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestQstick(t *testing.T) {
	type Data struct {
		Open   float64
		Close  float64
		Qstick float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/qstick.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	openings := helper.Map(inputs[0], func(d *Data) float64 { return d.Open })
	closings := helper.Map(inputs[1], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[2], func(d *Data) float64 { return d.Qstick })

	qstick := momentum.NewQstick[float64]()
	actual := qstick.Compute(openings, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, qstick.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
