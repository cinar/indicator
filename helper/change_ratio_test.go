// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestChangeRatio(t *testing.T) {
	input := helper.SliceToChan([]float64{1, 2, 5, 5, 8, 2, 1, 1, 3, 4})
	expected := helper.SliceToChan([]float64{4, 1.5, 0.6, -0.6, -0.875, -0.5, 2, 3})

	actual := helper.ChangeRatio(input, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
