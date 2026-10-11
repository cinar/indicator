// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestBuffered(_ *testing.T) {
	c := make(chan int, 1)
	b := helper.Buffered(c, 4)

	c <- 1
	c <- 2
	c <- 3
	c <- 4

	close(c)

	helper.Drain(b)
}

func TestBufferedNegativeSize(t *testing.T) {
	input := helper.SliceToChan([]int{2, 4, 6, 8})
	expected := helper.SliceToChan([]int{2, 4, 6, 8})

	actual := helper.Buffered(input, -1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
