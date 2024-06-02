// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestMultiply(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 4, 2, 4, 2})
	bc := helper.SliceToChan([]int{2, 1, 3, 2, 5})

	expected := helper.SliceToChan([]int{2, 4, 6, 8, 10})

	actual := helper.Multiply(ac, bc)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
