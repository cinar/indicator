// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestPow(t *testing.T) {
	input := []int{2, 3, 5, 10}
	expected := []int{4, 9, 25, 100}

	actual := helper.ChanToSlice(helper.Pow(helper.SliceToChan(input), 2))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
