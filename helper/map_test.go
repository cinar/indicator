// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestMap(t *testing.T) {
	type Row struct {
		High int
		Low  int
	}

	input := []Row{
		{High: 10, Low: 5},
		{High: 20, Low: 15},
	}

	expected := []int{5, 15}

	actual := helper.ChanToSlice(helper.Map(helper.SliceToChan(input), func(r Row) int {
		return r.Low
	}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
