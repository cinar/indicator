// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestPo(t *testing.T) {
	type Data struct {
		High  float64
		Low   float64
		Close float64
		Po    float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/po.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[3], func(d *Data) float64 { return d.Po })

	po := volatility.NewPoWithPeriod[float64](50)
	actual := po.Compute(highs, lows, closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, po.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPoFlatMarket(t *testing.T) {
	po := volatility.NewPo[float64]()
	length := po.IdlePeriod() + 20

	highs := make([]float64, length)
	lows := make([]float64, length)
	closings := make([]float64, length)

	for i := range highs {
		highs[i] = 100
		lows[i] = 100
		closings[i] = 100
	}

	actual := helper.ChanToSlice(po.Compute(
		helper.SliceToChan(highs),
		helper.SliceToChan(lows),
		helper.SliceToChan(closings),
	))

	if len(actual) == 0 {
		t.Fatal("expected at least one PO value")
	}

	for i, v := range actual {
		if v != 50 {
			t.Fatalf("expected PO of 50 for flat market at %d, got %v", i, v)
		}
	}
}

func TestPoString(t *testing.T) {
	expected := "PO(14)"
	actual := volatility.NewPo[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
