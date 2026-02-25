// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSeq(t *testing.T) {
	expected := helper.SliceToChan([]int{2, 3, 4, 5})
	actual := helper.Seq(2, 6, 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeqIncrement(t *testing.T) {
	expected := helper.SliceToChan([]int{2, 4})
	actual := helper.Seq(2, 6, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeqEmpty(t *testing.T) {
	actual := helper.Seq(6, 2, 1)
	res := helper.ChanToSlice(actual)
	if len(res) != 0 {
		t.Fatalf("expected empty, got %v", res)
	}
}
