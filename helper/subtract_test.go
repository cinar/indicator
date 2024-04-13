// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSubstract(t *testing.T) {
	ac := helper.SliceToChan([]int{2, 4, 6, 8, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	actual := helper.Subtract(ac, bc)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
