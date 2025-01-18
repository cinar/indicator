// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volume"
)

func TestFi(t *testing.T) {
	type FiData struct {
		Close  float64
		Volume int64
		Fi     float64
	}

	input, err := helper.ReadFromCsvFile[FiData]("testdata/fi.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	closings := helper.Map(inputs[0], func(m *FiData) float64 { return m.Close })
	volumes := helper.Map(inputs[1], func(m *FiData) float64 { return float64(m.Volume) })
	expected := helper.Map(inputs[2], func(m *FiData) float64 { return m.Fi })

	fi := volume.NewFi[float64]()
	actual := helper.RoundDigits(fi.Compute(closings, volumes), 2)
	expected = helper.Skip(expected, fi.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
