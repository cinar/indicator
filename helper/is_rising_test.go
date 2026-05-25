// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestIsRising(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4})
	expected := helper.SliceToChan([]int{1, 1, 1, 0, 0, 0, 1, 1})

	actual := helper.IsRising(input, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsRisingDefaultPeriod(t *testing.T) {
	input := helper.SliceToChan([]int{1, 3, 2, 5, 4})
	expected := helper.SliceToChan([]int{1, 0, 1, 0})

	actual := helper.IsRising(input, 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsRisingAllRising(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	expected := helper.SliceToChan([]int{1, 1, 1, 1})

	actual := helper.IsRising(input, 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsRisingNoneRising(t *testing.T) {
	input := helper.SliceToChan([]int{5, 4, 3, 2, 1})
	expected := helper.SliceToChan([]int{0, 0, 0, 0})

	actual := helper.IsRising(input, 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsRisingEqual(t *testing.T) {
	input := helper.SliceToChan([]int{3, 3, 3, 3})
	expected := helper.SliceToChan([]int{0, 0, 0})

	actual := helper.IsRising(input, 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
