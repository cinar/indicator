// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"testing"
)

func TestLowest(t *testing.T) {
	input := SliceToChan([]int{48, 52, 50, 49, 10})
	expected := SliceToChan([]int{48, 48, 48, 49, 10})
	window := 3
	actual := Lowest(input, window)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
