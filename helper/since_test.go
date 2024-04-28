// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSince(t *testing.T) {
	input := helper.SliceToChan([]int{1, 1, 2, 2, 2, 1, 2, 3, 3, 4})
	expected := helper.SliceToChan([]int{0, 1, 0, 1, 2, 0, 0, 0, 1, 0})

	actual := helper.Since(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
