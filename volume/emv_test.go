// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestEmv(t *testing.T) {
	type EmvData struct {
		High   float64
		Low    float64
		Volume int64
		Emv    float64
	}

	input, err := helper.ReadFromCsvFile[EmvData]("testdata/emv.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	highs := helper.Map(inputs[0], func(m *EmvData) float64 { return m.High })
	lows := helper.Map(inputs[1], func(m *EmvData) float64 { return m.Low })
	volumes := helper.Map(inputs[2], func(m *EmvData) float64 { return float64(m.Volume) })
	expected := helper.Map(inputs[3], func(m *EmvData) float64 { return m.Emv })

	emv := volume.NewEmv[float64]()
	actual := helper.RoundDigits(emv.Compute(highs, lows, volumes), 2)
	expected = helper.Skip(expected, emv.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
