// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSqrt(t *testing.T) {
	input := []int{9, 81, 16, 100}
	expected := []int{3, 9, 4, 10}

	actual := helper.ChanToSlice(helper.Sqrt(helper.SliceToChan(input)))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
