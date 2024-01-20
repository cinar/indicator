// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestKeepPositives(t *testing.T) {
	input := []int{-10, 20, 4, -5}
	expected := []int{0, 20, 4, 0}

	actual := helper.ChanToSlice(helper.KeepPositives(helper.SliceToChan(input)))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
