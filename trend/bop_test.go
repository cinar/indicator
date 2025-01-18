// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestBop(t *testing.T) {
	type BopData struct {
		Open  float64
		High  float64
		Low   float64
		Close float64
		Bop   float64
	}

	input, err := helper.ReadFromCsvFile[BopData]("testdata/bop.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	opening := helper.Map(inputs[0], func(row *BopData) float64 { return row.Open })
	high := helper.Map(inputs[1], func(row *BopData) float64 { return row.High })
	low := helper.Map(inputs[2], func(row *BopData) float64 { return row.Low })
	closing := helper.Map(inputs[3], func(row *BopData) float64 { return row.Close })
	expected := helper.Map(inputs[4], func(row *BopData) float64 { return row.Bop })

	bop := trend.NewBop[float64]()
	actual := bop.Compute(opening, high, low, closing)

	actual = helper.RoundDigits(actual, 0)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
