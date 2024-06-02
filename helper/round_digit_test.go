// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRoundDigit(t *testing.T) {
	input := 10.1234
	expected := 10.12

	actual := helper.RoundDigit(input, 2)

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
