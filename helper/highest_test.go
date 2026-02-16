// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"testing"
)

func TestHighest(t *testing.T) {
	input := SliceToChan([]int{48, 52, 50, 49, 10})
	expected := SliceToChan([]int{48, 52, 52, 52, 50})
	window := 3
	actual := Highest(input, window)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
