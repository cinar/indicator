// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"context"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestInternalBarStrength(t *testing.T) {
	type Data struct {
		High  float64 `header:"High"`
		Low   float64 `header:"Low"`
		Close float64 `header:"Close"`
		Ibs   float64 `header:"Ibs"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ibs.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Ibs })

	ibs := momentum.NewInternalBarStrength[float64]()

	if ibs.IdlePeriod() != 0 {
		t.Fatalf("expected IdlePeriod to be 0, got %d", ibs.IdlePeriod())
	}

	if ibs.String() != "IBS" {
		t.Fatalf("expected String to be IBS, got %s", ibs.String())
	}

	actual := ibs.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.RoundDigits(expected, 2)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInternalBarStrengthZeroDenom(t *testing.T) {
	highs := make(chan float64, 1)
	lows := make(chan float64, 1)
	closings := make(chan float64, 1)

	highs <- 100.0
	lows <- 100.0
	closings <- 100.0

	close(highs)
	close(lows)
	close(closings)

	ibs := momentum.NewInternalBarStrength[float64]()
	actual := ibs.Compute(highs, lows, closings)

	val := <-actual
	if val != 0.5 {
		t.Fatalf("expected neutral 0.5 when high == low, got %f", val)
	}
}

func TestInternalBarStrengthFlatBarsFloat32(t *testing.T) {
	ibs := momentum.NewInternalBarStrength[float32]()
	actual := ibs.ComputeWithContext(context.Background(),
		helper.SliceToChan([]float32{100, 102, 100, 104}),
		helper.SliceToChan([]float32{100, 100, 100, 100}),
		helper.SliceToChan([]float32{100, 100, 100, 104}))
	expected := helper.SliceToChan([]float32{0.5, 0, 0.5, 1})
	if err := helper.CheckEquals(actual, expected); err != nil {
		t.Fatal(err)
	}
}
