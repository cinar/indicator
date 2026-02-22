// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestCfo(t *testing.T) {
	type CfoData struct {
		Close float64
		Cfo   float64
	}

	input, err := helper.ReadFromCsvFile[CfoData]("testdata/cfo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(a *CfoData) float64 { return a.Close })
	expected := helper.Map(inputs[1], func(a *CfoData) float64 { return a.Cfo })

	cfo := trend.NewCfoWithPeriod[float64](5)
	actual := helper.RoundDigits(cfo.Compute(closing), 2)
	expected = helper.RoundDigits(helper.Skip(expected, cfo.IdlePeriod()), 2)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
