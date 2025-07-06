// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestWeightedClose(t *testing.T) {
	type Data struct {
		High          float64
		Low           float64
		Close         float64
		WeightedClose float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/weighted_close.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Close })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.WeightedClose })

	weightedClose := trend.NewWeightedClose[float64]()

	actual := weightedClose.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, weightedClose.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWeightedCloseString(t *testing.T) {
	expected := "Weighted Close"
	actual := trend.NewWeightedClose[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
