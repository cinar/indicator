// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

// TestWma tests the WMA indicator with 2 different periods, WMA(3) and WMA(5).
func TestWma(t *testing.T) {
	type Data struct {
		Close float64
		Wma3  float64
		Wma5  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/wma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	closing := helper.Duplicate(helper.Map(inputs[0], func(d *Data) float64 { return d.Close }), 2)
	expectedWma3 := helper.Map(inputs[1], func(d *Data) float64 { return d.Wma3 })
	expectedWma5 := helper.Map(inputs[2], func(d *Data) float64 { return d.Wma5 })

	wma3 := trend.NewWmaWith[float64](3)
	wma5 := trend.NewWmaWith[float64](5)

	actualWma3 := wma3.Compute(closing[0])
	actualWma5 := wma5.Compute(closing[1])

	actualWma3 = helper.RoundDigits(actualWma3, 3)
	actualWma5 = helper.RoundDigits(actualWma5, 3)

	expectedWma3 = helper.Skip(expectedWma3, wma3.IdlePeriod())
	expectedWma5 = helper.Skip(expectedWma5, wma5.IdlePeriod())

	err = helper.CheckEquals(actualWma3, expectedWma3, actualWma5, expectedWma5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWmaString(t *testing.T) {
	expected := "WMA(10)"
	actual := trend.NewWmaWith[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
