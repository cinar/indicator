// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestBollingerBands(t *testing.T) {
	type Data struct {
		Close  float64
		Upper  float64
		Middle float64
		Lower  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/bollinger_bands.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	upper := helper.Map(inputs[1], func(d *Data) float64 { return d.Upper })
	middle := helper.Map(inputs[2], func(d *Data) float64 { return d.Middle })
	lower := helper.Map(inputs[3], func(d *Data) float64 { return d.Lower })

	bb := volatility.NewBollingerBands[float64]()
	actualUpper, actualMiddle, actualLower := bb.Compute(closings)
	actualUpper = helper.RoundDigits(actualUpper, 2)
	actualMiddle = helper.RoundDigits(actualMiddle, 2)
	actualLower = helper.RoundDigits(actualLower, 2)

	upper = helper.Skip(upper, bb.IdlePeriod())
	middle = helper.Skip(middle, bb.IdlePeriod())
	lower = helper.Skip(lower, bb.IdlePeriod())

	err = helper.CheckEquals(actualUpper, upper, actualMiddle, middle, actualLower, lower)
	if err != nil {
		t.Fatal(err)
	}
}
