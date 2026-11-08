// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMovingMax(t *testing.T) {
	input := helper.SliceToChan([]int{-10, 20, -4, -5, 1, 5, 8, 10, -20, 4})
	expected := helper.SliceToChan([]int{20, 20, 5, 8, 10, 10, 10})

	movingMax := trend.NewMovingMaxWithPeriod[int](4)
	actual := movingMax.Compute(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// TestMovingMaxZeroInFirstWindow checks that a 0 in the first window counts
// toward the maximum. The shifted stream used to remove values starts with
// zero fill values, which used to remove the real 0.
func TestMovingMaxZeroInFirstWindow(t *testing.T) {
	input := helper.SliceToChan([]float64{-3, 0, -1, -2, -4})
	expected := helper.SliceToChan([]float64{0, 0, -1})

	movingMax := trend.NewMovingMaxWithPeriod[float64](3)
	actual := movingMax.Compute(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMovingMaxString(t *testing.T) {
	expected := "MOVINGMAX(4)"
	actual := trend.NewMovingMaxWithPeriod[int](4).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
