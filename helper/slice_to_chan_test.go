// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSliceToChan(t *testing.T) {
	expected := []int{2, 4, 6, 8}
	actual := helper.ChanToSlice(helper.SliceToChan(expected))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
