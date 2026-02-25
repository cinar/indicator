// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestElderRay(t *testing.T) {
	type Data struct {
		High      float64 `header:"High"`
		Low       float64 `header:"Low"`
		Close     float64 `header:"Close"`
		BullPower float64 `header:"BullPower"`
		BearPower float64 `header:"BearPower"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/elder_ray.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	dataForCompute := helper.Duplicate(inputs[0], 3)
	highs := helper.Map(dataForCompute[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(dataForCompute[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(dataForCompute[2], func(d *Data) float64 { return d.Close })

	er := momentum.NewElderRayWithPeriod[float64](3)

	expectedString := "Elder-Ray Index(3)"
	if er.String() != expectedString {
		t.Fatalf("String: actual %s expected %s", er.String(), expectedString)
	}

	actualBullPower, actualBearPower := er.Compute(highs, lows, closings)

	actualBullPower = helper.RoundDigits(actualBullPower, 2)
	actualBearPower = helper.RoundDigits(actualBearPower, 2)

	inputs[1] = helper.Skip(inputs[1], er.IdlePeriod())

	for data := range inputs[1] {
		actualBull := <-actualBullPower
		actualBear := <-actualBearPower

		if actualBull != data.BullPower {
			t.Fatalf("Bull Power: actual %v expected %v", actualBull, data.BullPower)
		}

		if actualBear != data.BearPower {
			t.Fatalf("Bear Power: actual %v expected %v", actualBear, data.BearPower)
		}
	}
}
