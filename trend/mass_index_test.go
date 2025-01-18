// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMassIndex(t *testing.T) {
	type Data struct {
		Open      float64
		Close     float64
		MassIndex float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/mass_index.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	openings := helper.Map(inputs[0], func(d *Data) float64 { return d.Open })
	closings := helper.Map(inputs[1], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[2], func(d *Data) float64 { return d.MassIndex })

	mi := trend.NewMassIndex[float64]()

	actual := mi.Compute(openings, closings)
	actual = helper.RoundDigits(actual, 2)
	actual = helper.Shift(actual, mi.IdlePeriod(), 0)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
