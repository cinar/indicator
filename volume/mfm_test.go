// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestMfm(t *testing.T) {
	type MfmData struct {
		High  float64
		Low   float64
		Close float64
		Mfm   float64
	}

	input, err := helper.ReadFromCsvFile[MfmData]("testdata/mfm.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(m *MfmData) float64 { return m.High })
	lows := helper.Map(inputs[1], func(m *MfmData) float64 { return m.Low })
	closings := helper.Map(inputs[2], func(m *MfmData) float64 { return m.Close })
	expected := helper.Map(inputs[3], func(m *MfmData) float64 { return m.Mfm })

	mfm := volume.NewMfm[float64]()
	actual := helper.RoundDigits(mfm.Compute(highs, lows, closings), 2)
	expected = helper.Skip(expected, mfm.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
