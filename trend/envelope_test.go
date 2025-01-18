// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestEnvelopeWithSma(t *testing.T) {
	type Data struct {
		Close  float64
		Upper  float64
		Middle float64
		Lower  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/envelope_sma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expectedUpper := helper.Map(inputs[1], func(d *Data) float64 { return d.Upper })
	expectedMiddle := helper.Map(inputs[2], func(d *Data) float64 { return d.Middle })
	expectedLower := helper.Map(inputs[3], func(d *Data) float64 { return d.Lower })

	envelope := trend.NewEnvelopeWithSma[float64]()
	actualUpper, actualMiddle, actualLower := envelope.Compute(closing)

	actualUpper = helper.RoundDigits(actualUpper, 2)
	actualMiddle = helper.RoundDigits(actualMiddle, 2)
	actualLower = helper.RoundDigits(actualLower, 2)

	expectedUpper = helper.Skip(expectedUpper, envelope.IdlePeriod())
	expectedMiddle = helper.Skip(expectedMiddle, envelope.IdlePeriod())
	expectedLower = helper.Skip(expectedLower, envelope.IdlePeriod())

	err = helper.CheckEquals(actualUpper, expectedUpper, actualMiddle, expectedMiddle, actualLower, expectedLower)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnvelopeWithEma(t *testing.T) {
	type Data struct {
		Close  float64
		Upper  float64
		Middle float64
		Lower  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/envelope_ema.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expectedUpper := helper.Map(inputs[1], func(d *Data) float64 { return d.Upper })
	expectedMiddle := helper.Map(inputs[2], func(d *Data) float64 { return d.Middle })
	expectedLower := helper.Map(inputs[3], func(d *Data) float64 { return d.Lower })

	envelope := trend.NewEnvelopeWithEma[float64]()
	actualUpper, actualMiddle, actualLower := envelope.Compute(closing)

	actualUpper = helper.RoundDigits(actualUpper, 2)
	actualMiddle = helper.RoundDigits(actualMiddle, 2)
	actualLower = helper.RoundDigits(actualLower, 2)

	expectedUpper = helper.Skip(expectedUpper, envelope.IdlePeriod())
	expectedMiddle = helper.Skip(expectedMiddle, envelope.IdlePeriod())
	expectedLower = helper.Skip(expectedLower, envelope.IdlePeriod())

	err = helper.CheckEquals(actualUpper, expectedUpper, actualMiddle, expectedMiddle, actualLower, expectedLower)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnvelopeString(t *testing.T) {
	expected := "Envelope(SMA(1),2)"

	envelope := trend.NewEnvelope(trend.NewSmaWithPeriod[float64](1), 2)
	actual := envelope.String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
