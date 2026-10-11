// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestDonchianChannel(t *testing.T) {
	type Data struct {
		High   float64
		Low    float64
		Upper  float64
		Middle float64
		Lower  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/donchian_channel.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	expectedUpper := helper.Map(inputs[2], func(d *Data) float64 { return d.Upper })
	expectedMiddle := helper.Map(inputs[3], func(d *Data) float64 { return d.Middle })
	expectedLower := helper.Map(inputs[4], func(d *Data) float64 { return d.Lower })

	dc := volatility.NewDonchianChannel[float64]()
	actualUpper, actualMiddle, actualLower := dc.Compute(highs, lows)
	actualUpper = helper.RoundDigits(actualUpper, 2)
	actualMiddle = helper.RoundDigits(actualMiddle, 2)
	actualLower = helper.RoundDigits(actualLower, 2)

	expectedUpper = helper.Skip(expectedUpper, dc.IdlePeriod())
	expectedMiddle = helper.Skip(expectedMiddle, dc.IdlePeriod())
	expectedLower = helper.Skip(expectedLower, dc.IdlePeriod())

	err = helper.CheckEquals(actualUpper, expectedUpper, actualMiddle, expectedMiddle, actualLower, expectedLower)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDonchianChannelString(t *testing.T) {
	expected := "DC(20)"
	actual := volatility.NewDonchianChannel[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
