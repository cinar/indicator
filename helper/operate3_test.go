// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestOperate3(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := []int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30}

	actual := helper.ChanToSlice(helper.Operate3(ac, bc, cc, func(a, b, c int) int {
		return a + b + c
	}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestOperate3FirstEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := []int{3, 6, 9, 12, 15, 18, 21, 24}

	actual := helper.ChanToSlice(helper.Operate3(ac, bc, cc, func(a, b, c int) int {
		return a + b + c
	}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestOperate3SecondEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := []int{3, 6, 9, 12, 15, 18, 21, 24}

	actual := helper.ChanToSlice(helper.Operate3(ac, bc, cc, func(a, b, c int) int {
		return a + b + c
	}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestOperate3ThirdEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8})

	expected := []int{3, 6, 9, 12, 15, 18, 21, 24}

	actual := helper.ChanToSlice(helper.Operate3(ac, bc, cc, func(a, b, c int) int {
		return a + b + c
	}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
