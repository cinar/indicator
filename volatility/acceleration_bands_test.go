// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"sync"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestAccelerationBands(t *testing.T) {
	type Data struct {
		High   float64
		Low    float64
		Close  float64
		Upper  float64
		Middle float64
		Lower  float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/acceleration_bands.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 6)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	upper := helper.Map(inputs[3], func(d *Data) float64 { return d.Upper })
	middle := helper.Map(inputs[4], func(d *Data) float64 { return d.Middle })
	lower := helper.Map(inputs[5], func(d *Data) float64 { return d.Lower })

	ab := volatility.NewAccelerationBands[float64]()
	actualUpper, actualMiddle, actualLower := ab.Compute(highs, lows, closings)
	actualUpper = helper.RoundDigits(actualUpper, 2)
	actualMiddle = helper.RoundDigits(actualMiddle, 2)
	actualLower = helper.RoundDigits(actualLower, 2)

	upper = helper.Skip(upper, ab.IdlePeriod())
	middle = helper.Skip(middle, ab.IdlePeriod())
	lower = helper.Skip(lower, ab.IdlePeriod())

	err = helper.CheckEquals(actualUpper, upper, actualMiddle, middle, actualLower, lower)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccelerationBandsZeroPricedInstrument(t *testing.T) {
	period := volatility.DefaultAccelerationBandsPeriod
	length := period + 5

	highs := make([]float64, length)
	lows := make([]float64, length)
	closings := make([]float64, length)

	ab := volatility.NewAccelerationBands[float64]()
	actualUpper, actualMiddle, actualLower := ab.Compute(
		helper.SliceToChan(highs),
		helper.SliceToChan(lows),
		helper.SliceToChan(closings),
	)

	// The three output channels share upstream duplicated/fanned-out
	// channels internally, so they must be drained concurrently rather
	// than one at a time, or the pipeline deadlocks.
	var upperValues, middleValues, lowerValues []float64

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		upperValues = helper.ChanToSlice(actualUpper)
	}()

	go func() {
		defer wg.Done()
		middleValues = helper.ChanToSlice(actualMiddle)
	}()

	go func() {
		defer wg.Done()
		lowerValues = helper.ChanToSlice(actualLower)
	}()

	wg.Wait()

	if len(upperValues) == 0 {
		t.Fatal("expected at least one value")
	}

	for i, v := range upperValues {
		if v != 0 {
			t.Fatalf("expected upper band 0 for zero-priced instrument at %d, got %v", i, v)
		}
	}

	for i, v := range middleValues {
		if v != 0 {
			t.Fatalf("expected middle band 0 for zero-priced instrument at %d, got %v", i, v)
		}
	}

	for i, v := range lowerValues {
		if v != 0 {
			t.Fatalf("expected lower band 0 for zero-priced instrument at %d, got %v", i, v)
		}
	}
}

func TestAccelerationBandsString(t *testing.T) {
	expected := "ACCELBANDS(20)"
	actual := volatility.NewAccelerationBands[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
