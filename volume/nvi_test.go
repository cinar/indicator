// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestNvi(t *testing.T) {
	type NviData struct {
		Close  float64
		Volume int64
		Nvi    float64
	}

	input, err := helper.ReadFromCsvFile[NviData]("testdata/nvi.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	closings := helper.Map(inputs[0], func(m *NviData) float64 { return m.Close })
	volumes := helper.Map(inputs[1], func(m *NviData) float64 { return float64(m.Volume) })
	expected := helper.Map(inputs[2], func(m *NviData) float64 { return m.Nvi })

	nvi := volume.NewNvi[float64]()
	actual := helper.RoundDigits(nvi.Compute(closings, volumes), 2)
	expected = helper.Skip(expected, nvi.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
