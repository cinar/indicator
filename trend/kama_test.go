// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestKama(t *testing.T) {
	type Data struct {
		Close float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/kama.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 1)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	t.Fatal(helper.ChanToSlice(closing))
}
