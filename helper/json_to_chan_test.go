// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"strings"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestJSONToChan(t *testing.T) {
	expected := helper.SliceToChan([]int{2, 4, 6, 8})
	input := "[2, 4, 6, 8]"

	actual := helper.JSONToChan[int](strings.NewReader(input))

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
