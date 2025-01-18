// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestChaikinOscillator(t *testing.T) {
	type Data struct {
		High   float64
		Low    float64
		Close  float64
		Volume int64
		Ad     float64
		Co     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/chaikin_oscillator.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 6)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	volumes := helper.Map(inputs[3], func(d *Data) float64 { return float64(d.Volume) })
	expectedAd := helper.Map(inputs[4], func(d *Data) float64 { return d.Ad })
	expectedCo := helper.Map(inputs[5], func(d *Data) float64 { return d.Co })

	co := momentum.NewChaikinOscillator[float64]()
	actualCo, actualAd := co.Compute(highs, lows, closings, volumes)
	actualCo = helper.RoundDigits(actualCo, 2)
	actualAd = helper.RoundDigits(actualAd, 2)

	expectedAd = helper.Skip(expectedAd, co.IdlePeriod())
	expectedCo = helper.Skip(expectedCo, co.IdlePeriod())

	err = helper.CheckEquals(actualAd, expectedAd, actualCo, expectedCo)
	if err != nil {
		t.Fatal(err)
	}
}
