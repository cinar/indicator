// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestAbs(t *testing.T) {
	input := helper.SliceToChan([]int{-10, 20, -4, -5})
	expected := helper.SliceToChan([]int{10, 20, 4, 5})

	actual := helper.Abs(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAbsInt64PrecisionBeyondFloat64Range(t *testing.T) {
	// Beyond float64's 53-bit exact-integer range (2^53). A lossy
	// float64 round-trip would not return this exact value.
	var n int64 = (1 << 62) + 1234567

	input := helper.SliceToChan([]int64{-n})
	expected := helper.SliceToChan([]int64{n})

	actual := helper.Abs(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAbsMinInt8Overflow(t *testing.T) {
	// The true absolute value of math.MinInt8 (128) does not fit in
	// an int8 (range [-128, 127]). This documents the inherent,
	// unfixable overflow behavior at the minimum value of a signed
	// integer type: it wraps back to itself.
	input := helper.SliceToChan([]int8{math.MinInt8})
	expected := helper.SliceToChan([]int8{math.MinInt8})

	actual := helper.Abs(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
