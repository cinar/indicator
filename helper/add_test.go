// Package helper contains the helper_test functions test.
//
// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2
package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestAdd(t *testing.T) {
	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	actual := helper.ChanToSlice(helper.Add(ac, bc))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
