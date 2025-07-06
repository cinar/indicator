// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestPpo(t *testing.T) {
	type Data struct {
		Close     float64
		Ppo       float64
		Signal    float64
		Histogram float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ppo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expectedPpo := helper.Map(inputs[1], func(d *Data) float64 { return d.Ppo })
	expectedSignal := helper.Map(inputs[2], func(d *Data) float64 { return d.Signal })
	expectedHistogram := helper.Map(inputs[3], func(d *Data) float64 { return d.Histogram })

	ppo := momentum.NewPpo[float64]()
	actualPpo, actualSignal, actualHistogram := ppo.Compute(closings)
	actualPpo = helper.RoundDigits(actualPpo, 2)
	actualSignal = helper.RoundDigits(actualSignal, 2)
	actualHistogram = helper.RoundDigits(actualHistogram, 2)

	expectedPpo = helper.Skip(expectedPpo, ppo.IdlePeriod())
	expectedSignal = helper.Skip(expectedSignal, ppo.IdlePeriod())
	expectedHistogram = helper.Skip(expectedHistogram, ppo.IdlePeriod())

	err = helper.CheckEquals(actualPpo, expectedPpo, actualSignal, expectedSignal, actualHistogram, expectedHistogram)
	if err != nil {
		t.Fatal(err)
	}
}
