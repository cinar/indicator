// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSync(t *testing.T) {
	input1 := helper.Skip(helper.SliceToChan([]int{0, 0, 0, 0, 1, 2, 3, 4}), 4)
	input2 := helper.Skip(helper.SliceToChan([]int{0, 0, 1, 2, 3, 4, 5, 6}), 2)
	input3 := helper.Skip(helper.SliceToChan([]int{0, 0, 0, 1, 2, 3, 4, 5}), 3)

	commonPeriod := helper.CommonPeriod(4, 2, 3)

	actual1 := helper.SyncPeriod(commonPeriod, 4, input1)
	expected1 := helper.SliceToChan([]int{1, 2, 3, 4})

	actual2 := helper.SyncPeriod(commonPeriod, 2, input2)
	expected2 := helper.SliceToChan([]int{3, 4, 5, 6})

	actual3 := helper.SyncPeriod(commonPeriod, 3, input3)
	expected3 := helper.SliceToChan([]int{2, 3, 4, 5})

	err := helper.CheckEquals(actual1, expected1, actual2, expected2, actual3, expected3)
	if err != nil {
		t.Fatal(err)
	}
}
