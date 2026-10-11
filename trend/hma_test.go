// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"fmt"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestHma(t *testing.T) {
	type Data struct {
		Close float64
		Hma   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/hma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Hma })

	hma := trend.NewHmaWithPeriod[float64](3)

	actual := hma.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, hma.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewHma(t *testing.T) {
	hma := trend.NewHma[float64]()
	expected := fmt.Sprintf("HMA(%d)", trend.DefaultHmaPeriod)
	if hma.String() != expected {
		t.Fatalf("actual %v expected %v", hma.String(), expected)
	}
}

func TestHmaString(t *testing.T) {
	expected := "HMA(10)"
	actual := trend.NewHmaWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
