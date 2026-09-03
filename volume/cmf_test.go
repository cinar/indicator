// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestCmf(t *testing.T) {
	type CmfData struct {
		High   float64
		Low    float64
		Close  float64
		Volume int64
		Cmf    float64
	}

	input, err := helper.ReadFromCsvFile[CmfData]("testdata/cmf.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	highs := helper.Map(inputs[0], func(m *CmfData) float64 { return m.High })
	lows := helper.Map(inputs[1], func(m *CmfData) float64 { return m.Low })
	closings := helper.Map(inputs[2], func(m *CmfData) float64 { return m.Close })
	volumes := helper.Map(inputs[3], func(m *CmfData) float64 { return float64(m.Volume) })
	expected := helper.Map(inputs[4], func(m *CmfData) float64 { return m.Cmf })

	cmf := volume.NewCmf[float64]()
	actual := helper.RoundDigits(cmf.Compute(highs, lows, closings, volumes), 2)
	expected = helper.Skip(expected, cmf.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCmfZeroVolume(t *testing.T) {
	highs := helper.SliceToChan([]float64{11, 11, 11, 11})
	lows := helper.SliceToChan([]float64{9, 9, 9, 9})
	closings := helper.SliceToChan([]float64{10, 10, 10, 10})
	volumes := helper.SliceToChan([]float64{0, 0, 0, 0})

	cmf := volume.NewCmfWithPeriod[float64](2)
	actuals := cmf.Compute(highs, lows, closings, volumes)

	for actual := range actuals {
		if actual != 0 {
			t.Fatalf("expected 0 when window volume is 0, got %v", actual)
		}
	}
}

func TestCmfString(t *testing.T) {
	expected := "CMF(20)"
	actual := volume.NewCmf[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
