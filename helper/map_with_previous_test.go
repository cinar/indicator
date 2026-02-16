// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestMapWithPrevious(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 3, 4})
	expected := helper.SliceToChan([]int{1, 3, 6, 10})

	actual := helper.MapWithPrevious(input, func(p, c int) int {
		return p + c
	}, 0)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
