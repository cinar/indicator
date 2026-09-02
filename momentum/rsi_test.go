// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestRsi(t *testing.T) {
	type Data struct {
		Close float64
		Rsi   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/rsi.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expectedRsi := helper.Map(inputs[1], func(d *Data) float64 { return d.Rsi })

	rsi := momentum.NewRsi[float64]()
	actualRsi := rsi.Compute(closings)
	actualRsi = helper.RoundDigits(actualRsi, 2)

	expectedRsi = helper.Skip(expectedRsi, rsi.IdlePeriod())

	err = helper.CheckEquals(actualRsi, expectedRsi)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRsiFlatMarket(t *testing.T) {
	closings := make([]float64, momentum.DefaultRsiPeriod+20)
	for i := range closings {
		closings[i] = 100
	}

	rsi := momentum.NewRsi[float64]()
	actualRsi := helper.ChanToSlice(rsi.Compute(helper.SliceToChan(closings)))

	if len(actualRsi) == 0 {
		t.Fatal("expected at least one RSI value")
	}

	for _, v := range actualRsi {
		if v != 50 {
			t.Fatalf("expected RSI of 50 for flat market, got %v", v)
		}
	}
}

func TestRsiString(t *testing.T) {
	expected := "RSI(14)"
	actual := momentum.NewRsi[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
