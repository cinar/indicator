// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestVwap(t *testing.T) {
	type VwapData struct {
		Close  float64
		Volume int64
		Vwap   float64
	}

	input, err := helper.ReadFromCsvFile[VwapData]("testdata/vwap.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	closings := helper.Map(inputs[0], func(m *VwapData) float64 { return m.Close })
	volumes := helper.Map(inputs[1], func(m *VwapData) float64 { return float64(m.Volume) })
	expected := helper.Map(inputs[2], func(m *VwapData) float64 { return m.Vwap })

	vwap := volume.NewVwap[float64]()
	actual := helper.RoundDigits(vwap.Compute(closings, volumes), 2)
	expected = helper.Skip(expected, vwap.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVwapZeroVolume(t *testing.T) {
	closings := helper.SliceToChan([]float64{10, 20, 30})
	volumes := helper.SliceToChan([]float64{0, 0, 0})

	vwap := volume.NewVwapWithPeriod[float64](2)
	actuals := vwap.Compute(closings, volumes)

	for actual := range actuals {
		if actual != 0 {
			t.Fatalf("expected 0 before any window has volume, got %v", actual)
		}
	}
}

func TestVwapCarriesForwardOnZeroVolume(t *testing.T) {
	closings := helper.SliceToChan([]float64{10, 20, 30, 40})
	volumes := helper.SliceToChan([]float64{5, 5, 0, 0})

	vwap := volume.NewVwapWithPeriod[float64](2)
	actuals := helper.ChanToSlice(vwap.Compute(closings, volumes))

	// Windows: (10,20)@vol(5,5) -> 15; (20,30)@vol(5,0) -> 20;
	// (30,40)@vol(0,0) -> carries forward the last valid VWAP, 20.
	expected := []float64{15, 20, 20}

	if len(actuals) != len(expected) {
		t.Fatalf("expected %d values but got %d", len(expected), len(actuals))
	}

	for i, actual := range actuals {
		if actual != expected[i] {
			t.Fatalf("index %d: expected %v but got %v", i, expected[i], actual)
		}
	}
}

func TestVwapString(t *testing.T) {
	expected := "VWAP(14)"
	actual := volume.NewVwap[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
