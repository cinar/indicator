// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestEcho(t *testing.T) {
	input := helper.SliceToChan([]int{2, 4, 6, 8})
	expected := helper.SliceToChan([]int{2, 4, 6, 8, 6, 8, 6, 8, 6, 8, 6, 8})

	actual := helper.Echo(input, 2, 4)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
