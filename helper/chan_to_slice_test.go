// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestChanToSlice(t *testing.T) {
	input := []int{2, 4, 6, 8}
	expected := helper.SliceToChan(input)

	actual := make(chan int, len(input))
	for _, n := range input {
		actual <- n
	}
	close(actual)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
