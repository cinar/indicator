// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestDivideBy(t *testing.T) {
	input := []int{2, 4, 6, 8}
	expected := []int{1, 2, 3, 4}

	actual := helper.ChanToSlice(helper.DivideBy(helper.SliceToChan(input), 2))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
