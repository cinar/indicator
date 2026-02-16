// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestDivideBy(t *testing.T) {
	input := []int{2, 4, 6, 8}
	expected := helper.SliceToChan([]int{1, 2, 3, 4})

	actual := helper.DivideBy(helper.SliceToChan(input), 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
