// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestPipe(t *testing.T) {
	data := []int{2, 4, 6, 8}
	expected := helper.SliceToChan(data)

	input := helper.SliceToChan(data)
	actual := make(chan int)

	go helper.Pipe(input, actual)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
