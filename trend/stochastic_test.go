// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestStochastic(t *testing.T) {
	type Data struct {
		Value float64
		K     float64
		D     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	values := helper.Map(inputs[0], func(d *Data) float64 { return d.Value })
	expectedK := helper.Map(inputs[1], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[2], func(d *Data) float64 { return d.D })

	s := trend.NewStochastic[float64]()
	actualK, actualD := s.Compute(values)
	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	expectedK = helper.Skip(expectedK, s.IdlePeriod())
	expectedD = helper.Skip(expectedD, s.IdlePeriod())

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD)
	if err != nil {
		t.Fatal(err)
	}
}
