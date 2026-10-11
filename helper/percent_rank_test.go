// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestPercentRank(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	// period = 5, window size = period-1 = 4.
	// The first 4 values (10, 20, 30, 40) are buffered as predecessors.
	// Emits starting from the 5th item (50), which already has 4 predecessors.
	// Window: [10, 20, 30, 40], value: 50 -> 100
	// Window: [20, 30, 40, 50], value: 60 -> 100
	// etc. (6 values total, one for each of the 5th through 10th items)
	expected := helper.SliceToChan([]float64{100, 100, 100, 100, 100, 100})
	actual := helper.PercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentRankMixed(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 50, 30, 20, 40, 15})
	// period = 5, window size = period-1 = 4.
	// The first 4 values (10, 50, 30, 20) are buffered as predecessors.
	// Emits starting from the 5th item (40), which already has 4 predecessors.
	// Window: [10, 50, 30, 20], value 40. Less than 40: 10, 30, 20 -> 3. 3*100/4 = 75.
	// Window: [50, 30, 20, 40], value 15. Less than 15: none -> 0. 0*100/4 = 0.
	expected := helper.SliceToChan([]float64{75, 0})
	actual := helper.PercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentRankPeriodThree(t *testing.T) {
	input := helper.SliceToChan([]float64{1, 5, 3, 4, 2, 6, 7})
	// period = 3, window size = period-1 = 2.
	// The first 2 values (1, 5) are buffered as predecessors.
	// The 3rd value (3) already has exactly 2 predecessors, so it must be ranked
	// immediately, using the full window [1, 5] (not a window that has already
	// dropped one of the two available predecessors).
	// Window: [1, 5], value 3. Less than 3: 1 -> 1. 1*100/2 = 50.
	// Window: [5, 3], value 4. Less than 4: 3 -> 1. 1*100/2 = 50.
	// Window: [3, 4], value 2. Less than 2: none -> 0. 0*100/2 = 0.
	// Window: [4, 2], value 6. Less than 6: 4, 2 -> 2. 2*100/2 = 100.
	// Window: [2, 6], value 7. Less than 7: 2, 6 -> 2. 2*100/2 = 100.
	expected := helper.SliceToChan([]float64{50, 50, 0, 100, 100})
	actual := helper.PercentRank(input, 3)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSortedPercentRank(t *testing.T) {
	input := helper.SliceToChan([]float64{10, 50, 30, 20, 40, 15})
	// period = 5, window size = period-1 = 4.
	// The first 4 values (10, 50, 30, 20) are buffered as predecessors.
	// Emits starting from the 5th item (40), which already has 4 predecessors.
	// Window: [10, 50, 30, 20], sorted [10, 20, 30, 50], value 40 -> index 3. 3*100/4 = 75.
	// Window: [50, 30, 20, 40], sorted [20, 30, 40, 50], value 15 -> index 0. 0*100/4 = 0.
	expected := helper.SliceToChan([]float64{75, 0})
	actual := helper.SortedPercentRank(input, 5)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSortedPercentRankPeriodThree(t *testing.T) {
	input := helper.SliceToChan([]float64{1, 5, 3, 4, 2, 6, 7})
	// period = 3, window size = period-1 = 2.
	// Mirrors TestPercentRankPeriodThree: the 3rd value (3) already has exactly 2
	// predecessors and must be ranked immediately against the full window [1, 5].
	expected := helper.SliceToChan([]float64{50, 50, 0, 100, 100})
	actual := helper.SortedPercentRank(input, 3)

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
