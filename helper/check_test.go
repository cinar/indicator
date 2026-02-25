// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

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

func TestCheckEqualsMultiplePairs(t *testing.T) {
	a1 := helper.SliceToChan([]int{1, 2})
	b1 := helper.SliceToChan([]int{1, 2})
	a2 := helper.SliceToChan([]int{3, 4, 5})
	b2 := helper.SliceToChan([]int{3, 4, 5})

	err := helper.CheckEquals(a1, b1, a2, b2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckEqualsMultiplePairsFail(t *testing.T) {
	a1 := helper.SliceToChan([]int{1, 2})
	b1 := helper.SliceToChan([]int{1, 2})
	a2 := helper.SliceToChan([]int{3, 4, 5})
	b2 := helper.SliceToChan([]int{3, 4, 6})

	err := helper.CheckEquals(a1, b1, a2, b2)
	if err == nil {
		t.Fatal("expected error for mismatch in second pair")
	}
}
