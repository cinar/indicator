// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestSmma(t *testing.T) {
	type Data struct {
		Close float64
		Smma  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/smma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Smma })

	smma := trend.NewSmma[float64]()

	actual := smma.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, smma.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSmmaString(t *testing.T) {
	expected := "SMMA(10)"
	actual := trend.NewSmmaWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
