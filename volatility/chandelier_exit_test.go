// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestChandelierExit(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		Long  float64
		Short float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/chandelier_exit.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedLong := helper.Map(inputs[3], func(d *Data) float64 { return d.Long })
	expectedShort := helper.Map(inputs[4], func(d *Data) float64 { return d.Short })

	ce := volatility.NewChandelierExit[float64]()
	actualLong, actualShort := ce.Compute(highs, lows, closings)
	actualLong = helper.RoundDigits(actualLong, 2)
	actualShort = helper.RoundDigits(actualShort, 2)

	expectedLong = helper.Skip(expectedLong, ce.IdlePeriod())
	expectedShort = helper.Skip(expectedShort, ce.IdlePeriod())

	err = helper.CheckEquals(actualShort, expectedShort, actualLong, expectedLong)
	if err != nil {
		t.Fatal(err)
	}
}
