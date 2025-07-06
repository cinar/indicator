// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestIchimokuCloud(t *testing.T) {
	type Data struct {
		High           float64
		Low            float64
		Close          float64
		ConversionLine float64
		BaseLine       float64
		LeadingSpanA   float64
		LeadingSpanB   float64
		LaggingLine    float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ichimoku_cloud.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 8)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expectedConversionLine := helper.Map(inputs[3], func(d *Data) float64 { return d.ConversionLine })
	expectedBaseLine := helper.Map(inputs[4], func(d *Data) float64 { return d.BaseLine })
	expectedLeadingSpanA := helper.Map(inputs[5], func(d *Data) float64 { return d.LeadingSpanA })
	expectedLeadingSpanB := helper.Map(inputs[6], func(d *Data) float64 { return d.LeadingSpanB })
	expectedLaggingLine := helper.Map(inputs[7], func(d *Data) float64 { return d.LaggingLine })

	ic := momentum.NewIchimokuCloud[float64]()
	actualConversionLine, actualBaseLine, actualLeadingSpanA, actualLeadingSpanB, actualLaggingLine := ic.Compute(highs, lows, closings)

	actualConversionLine = helper.RoundDigits(actualConversionLine, 2)
	actualBaseLine = helper.RoundDigits(actualBaseLine, 2)
	actualLeadingSpanA = helper.RoundDigits(actualLeadingSpanA, 2)
	actualLeadingSpanB = helper.RoundDigits(actualLeadingSpanB, 2)
	actualLaggingLine = helper.RoundDigits(actualLaggingLine, 2)

	expectedConversionLine = helper.Skip(expectedConversionLine, ic.IdlePeriod())
	expectedBaseLine = helper.Skip(expectedBaseLine, ic.IdlePeriod())
	expectedLeadingSpanA = helper.Skip(expectedLeadingSpanA, ic.IdlePeriod())
	expectedLeadingSpanB = helper.Skip(expectedLeadingSpanB, ic.IdlePeriod())
	expectedLaggingLine = helper.Skip(expectedLaggingLine, ic.IdlePeriod())

	err = helper.CheckEquals(
		actualConversionLine, expectedConversionLine,
		actualBaseLine, expectedBaseLine,
		actualLeadingSpanA, expectedLeadingSpanA,
		actualLeadingSpanB, expectedLeadingSpanB,
		actualLaggingLine, expectedLaggingLine,
	)
	if err != nil {
		t.Fatal(err)
	}
}
