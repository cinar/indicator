// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
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

	expected := helper.SliceToChan([]int{5, 15})

	actual := helper.Map(helper.SliceToChan(input), func(r Row) int {
		return r.Low
	})

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
