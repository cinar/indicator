// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestChanToSlice(t *testing.T) {
	expected := []int{2, 4, 6, 8}

	c := make(chan int, len(expected))
	for _, n := range expected {
		c <- n
	}
	close(c)

	actual := helper.ChanToSlice(c)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
