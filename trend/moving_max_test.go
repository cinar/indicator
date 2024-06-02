// Copyright (c) 2021-2024 Onur Cinar.
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

	max := trend.NewMovingMaxWithPeriod[int](4)
	actual := max.Compute(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
