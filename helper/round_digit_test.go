// Copyright (c) 2021-2026 The Indicator Authors.
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

func TestRoundDigitIntegerIsNoOp(t *testing.T) {
	// Beyond float64's 53-bit exact-integer range (2^53). A lossy
	// float64 round-trip would not return this exact value.
	var input int64 = (1 << 62) + 1234567
	expected := input

	actual := helper.RoundDigit(input, 2)

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
