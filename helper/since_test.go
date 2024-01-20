// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSince(t *testing.T) {
	input := []int{1, 1, 2, 2, 2, 1, 2, 3, 3, 4}
	expected := []int{0, 1, 0, 1, 2, 0, 0, 0, 1, 0}

	actual := helper.ChanToSlice(helper.Since(helper.SliceToChan(input)))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
