// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestAbs(t *testing.T) {
	input := helper.SliceToChan([]int{-10, 20, -4, -5})
	expected := helper.SliceToChan([]int{10, 20, 4, 5})

	actual := helper.Abs(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
