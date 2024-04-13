// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestOperate(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := helper.SliceToChan([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})

	actual := helper.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperateFirstEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := helper.SliceToChan([]int{2, 4, 6, 8, 10, 12, 14, 16})

	actual := helper.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperateSecondEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8})

	expected := helper.SliceToChan([]int{2, 4, 6, 8, 10, 12, 14, 16})

	actual := helper.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
