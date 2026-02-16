// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestShift(t *testing.T) {
	input := helper.SliceToChan([]int{2, 4, 6, 8})
	expected := helper.SliceToChan([]int{0, 0, 0, 0, 2, 4, 6, 8})

	actual := helper.Shift(input, 4, 0)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
