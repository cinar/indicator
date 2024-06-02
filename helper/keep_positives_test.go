// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestKeepPositives(t *testing.T) {
	input := []int{-10, 20, 4, -5}
	expected := helper.SliceToChan([]int{0, 20, 4, 0})

	actual := helper.KeepPositives(helper.SliceToChan(input))

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
