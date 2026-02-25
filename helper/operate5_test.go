// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestOperate5(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	ec := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{5, 10, 15, 20, 25})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate5FirstEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	ec := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{5, 10, 15, 20})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate5SecondEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	ec := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{5, 10, 15, 20})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate5ThirdEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	ec := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{5, 10, 15, 20})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate5FourthEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4})
	ec := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{5, 10, 15, 20})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate5FifthEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	ec := helper.SliceToChan([]int{1, 2, 3, 4})

	expected := helper.SliceToChan([]int{5, 10, 15, 20})

	actual := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
