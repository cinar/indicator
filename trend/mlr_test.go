// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMlr(t *testing.T) {
	type Data struct {
		X float64
		Y float64
		R float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/mlr.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	x := helper.Map(inputs[0], func(d *Data) float64 { return d.X })
	y := helper.Map(inputs[1], func(d *Data) float64 { return d.Y })
	expected := helper.Map(inputs[2], func(d *Data) float64 { return d.R })

	mlr := trend.NewMlrWithPeriod[float64](4)

	actual := mlr.Compute(x, y)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, mlr.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
