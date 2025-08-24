// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestAtr(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		Atr   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/atr.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Atr })

	atr := volatility.NewAtrWithPeriod[float64](50)
	actual := atr.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, atr.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
