// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestPvo(t *testing.T) {
	type Data struct {
		Volume    float64
		Pvo       float64
		Signal    float64
		Histogram float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/pvo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Volume })
	expectedPvo := helper.Map(inputs[1], func(d *Data) float64 { return d.Pvo })
	expectedSignal := helper.Map(inputs[2], func(d *Data) float64 { return d.Signal })
	expectedHistogram := helper.Map(inputs[3], func(d *Data) float64 { return d.Histogram })

	pvo := momentum.NewPvo[float64]()
	actualPvo, actualSignal, actualHistogram := pvo.Compute(closings)
	actualPvo = helper.RoundDigits(actualPvo, 2)
	actualSignal = helper.RoundDigits(actualSignal, 2)
	actualHistogram = helper.RoundDigits(actualHistogram, 2)

	expectedPvo = helper.Skip(expectedPvo, pvo.IdlePeriod())
	expectedSignal = helper.Skip(expectedSignal, pvo.IdlePeriod())
	expectedHistogram = helper.Skip(expectedHistogram, pvo.IdlePeriod())

	err = helper.CheckEquals(actualPvo, expectedPvo, actualSignal, expectedSignal, actualHistogram, expectedHistogram)
	if err != nil {
		t.Fatal(err)
	}
}
