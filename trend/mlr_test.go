// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestNewMlr(t *testing.T) {
	mlr := trend.NewMlr[float64]()
	if mlr.Mls.Sum.Period != trend.DefaultMlrPeriod {
		t.Fatalf("actual %v expected %v", mlr.Mls.Sum.Period, trend.DefaultMlrPeriod)
	}
}

func TestNewMlrWithPeriod(t *testing.T) {
	mlr := trend.NewMlrWithPeriod[float64](10)
	if mlr.Mls.Sum.Period != 10 {
		t.Fatalf("actual %v expected %v", mlr.Mls.Sum.Period, 10)
	}
}

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

func TestMlrString(t *testing.T) {
	expected := "MLR(4)"
	actual := trend.NewMlrWithPeriod[float64](4).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
