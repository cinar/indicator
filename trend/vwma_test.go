// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestVwma(t *testing.T) {
	type Data struct {
		Close  float64
		Volume int64
		Vwma   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/vwma.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	volume := helper.Map(inputs[1], func(d *Data) float64 { return float64(d.Volume) })
	expected := helper.Map(inputs[2], func(d *Data) float64 { return d.Vwma })

	vwma := trend.NewVwma[float64]()

	actual := vwma.Compute(closing, volume)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, vwma.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVwmaNoVolume(t *testing.T) {
	vwma := trend.NewVwma[float64]()

	count := vwma.Period + 5
	closing := make([]float64, count)
	volume := make([]float64, count)

	for i := range count {
		closing[i] = 100
		volume[i] = 0
	}

	actual := helper.ChanToSlice(vwma.Compute(helper.SliceToChan(closing), helper.SliceToChan(volume)))

	if len(actual) == 0 {
		t.Fatal("expected at least one VWMA value")
	}

	for _, v := range actual {
		if v != 0 {
			t.Fatalf("expected VWMA of 0 for a window with no volume, got %v", v)
		}
	}
}

func TestVwmaString(t *testing.T) {
	expected := "VWMA(20)"
	actual := trend.NewVwma[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
