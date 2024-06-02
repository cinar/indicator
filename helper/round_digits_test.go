// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRoundDigits(t *testing.T) {
	input := helper.SliceToChan([]float64{10.1234, 5.678, 6.78, 8.91011})
	expected := helper.SliceToChan([]float64{10.12, 5.68, 6.78, 8.91})

	actual := helper.RoundDigits(input, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
