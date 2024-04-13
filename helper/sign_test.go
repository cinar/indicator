// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSign(t *testing.T) {
	input := helper.SliceToChan([]int{-10, 20, -4, 0})
	expected := helper.SliceToChan([]int{-1, 1, -1, 0})

	actual := helper.Sign(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
