// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestPipe(t *testing.T) {
	expected := []int{2, 4, 6, 8}
	input := helper.SliceToChan(expected)
	output := make(chan int)

	go helper.Pipe(input, output)
	actual := helper.ChanToSlice(output)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
