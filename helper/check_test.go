// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestCheckEqualsNotPairs(t *testing.T) {
	c := helper.SliceToChan([]int{1, 2, 3, 4})

	err := helper.CheckEquals(c)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEqualsNotEndedTheSame(t *testing.T) {
	a := helper.SliceToChan([]int{1, 2, 3, 4})
	b := helper.SliceToChan([]int{1, 2})

	err := helper.CheckEquals(a, b)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEqualsNotEquals(t *testing.T) {
	a := helper.SliceToChan([]int{1, 2, 3, 4})
	b := helper.SliceToChan([]int{1, 3, 3, 4})

	err := helper.CheckEquals(a, b)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEquals(t *testing.T) {
	a := helper.SliceToChan([]int{1, 2, 3, 4})
	b := helper.SliceToChan([]int{1, 2, 3, 4})

	err := helper.CheckEquals(a, b)
	if err != nil {
		t.Fatal(err)
	}
}
