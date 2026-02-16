// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestFirst(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := helper.SliceToChan([]int{1, 2, 3, 4})

	actual := helper.First(input, 4)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFirstLessValues(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2})
	expected := helper.SliceToChan([]int{1, 2})

	actual := helper.First(input, 4)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
