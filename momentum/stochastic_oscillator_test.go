// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestStochasticOscillator(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		K     float64
		D     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic_oscillator.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedK := helper.Map(inputs[3], func(d *Data) float64 { return d.K })
	expectedD := helper.Map(inputs[4], func(d *Data) float64 { return d.D })

	so := momentum.NewStochasticOscillator[float64]()
	actualK, actualD := so.Compute(highs, lows, closings)
	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	expectedK = helper.Skip(expectedK, so.IdlePeriod())
	expectedD = helper.Skip(expectedD, so.IdlePeriod())

	err = helper.CheckEquals(actualK, expectedK, actualD, expectedD)
	if err != nil {
		t.Fatal(err)
	}
}
