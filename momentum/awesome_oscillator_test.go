// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestAwesomeOscillator(t *testing.T) {
	type Data struct {
		High float64
		Low  float64
		Ao   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/awesome_oscillator.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	expected := helper.Map(inputs[2], func(d *Data) float64 { return d.Ao })

	ao := momentum.NewAwesomeOscillator[float64]()
	actual := ao.Compute(highs, lows)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, ao.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
