// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestPercentRank(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	// period = 5.
	// Emits from the 6th item (60).
	// Window: [20, 30, 40, 50], value: 60 -> 100
	// Window: [30, 40, 50, 60], value: 70 -> 100
	// etc. (5 values total)
	expected := helper.SliceToChan([]float64{100, 100, 100, 100, 100})
	actual := helper.PercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentRankMixed(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 50, 30, 20, 40, 15})
	// period = 5.
	// Emits from the 6th item (15).
	// Window: [50, 30, 20, 40], value 15. 0 less. 0*100/4 = 0.
	expected := helper.SliceToChan([]float64{0})
	actual := helper.PercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSortedPercentRank(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 50, 30, 20, 40, 15})
	// period = 5.
	// Emits from the 6th item (15).
	// Window: [50, 30, 20, 40], value 15. Sorted [20, 30, 40, 50]. Rank for 15 is index 0. 0*100/4 = 0.
	expected := helper.SliceToChan([]float64{0})
	actual := helper.SortedPercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentRankPeriodOne(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 20, 30})
	actual := helper.PercentRank(input, 1)
	// Should be empty
	res := helper.ChanToSlice(actual)
	if len(res) != 0 {
		t.Fatalf("expected empty, got %v", res)
	}
}

func TestSortedPercentRankPeriodOne(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 20, 30})
	actual := helper.SortedPercentRank(input, 1)
	// Should be empty
	res := helper.ChanToSlice(actual)
	if len(res) != 0 {
		t.Fatalf("expected empty, got %v", res)
	}
}
