// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestOperate4(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{4, 8, 12, 16, 20})

	actual := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
		return a + b + c + d
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate4FirstEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{4, 8, 12, 16})

	actual := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
		return a + b + c + d
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate4SecondEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{4, 8, 12, 16})

	actual := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
		return a + b + c + d
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate4ThirdEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4})
	dc := helper.SliceToChan([]int{1, 2, 3, 4, 5})

	expected := helper.SliceToChan([]int{4, 8, 12, 16})

	actual := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
		return a + b + c + d
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOperate4FourthEnds(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	cc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
	dc := helper.SliceToChan([]int{1, 2, 3, 4})

	expected := helper.SliceToChan([]int{4, 8, 12, 16})

	actual := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
		return a + b + c + d
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
