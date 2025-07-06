// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestKdj(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		K     float64
		D     float64
		J     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/kdj.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 6)
	high := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	low := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closing := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedK := helper.Map(inputs[3], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[4], func(d *Data) float64 { return d.D })
	expectedJ := helper.Map(inputs[5], func(d *Data) float64 { return d.J })

	kdj := trend.NewKdj[float64]()
	actualK, actualD, actualJ := kdj.Compute(high, low, closing)

	actualK = helper.RoundDigits(actualK, 2)
	actualK = helper.Shift(actualK, kdj.IdlePeriod(), 0)

	actualD = helper.RoundDigits(actualD, 2)
	actualD = helper.Shift(actualD, kdj.IdlePeriod(), 0)

	actualJ = helper.RoundDigits(actualJ, 2)
	actualJ = helper.Shift(actualJ, kdj.IdlePeriod(), 0)

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD, actualJ, expectedJ)
	if err != nil {
		t.Fatal(err)
	}
}
