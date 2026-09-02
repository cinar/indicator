// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestLcmWithTwoValues(t *testing.T) {
	actual := helper.Lcm(18, 32)
	expected := 288

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestLcmWithFourValues(t *testing.T) {
	actual := helper.Lcm(1, 2, 8, 6)
	expected := 24

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestLcmWithFiveValues(t *testing.T) {
	actual := helper.Lcm(2, 7, 3, 9, 8)
	expected := 504

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestLcmWithEmptyInput(t *testing.T) {
	actual := helper.Lcm()
	expected := 0

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestLcmWithZeroInInput(t *testing.T) {
	actual := helper.Lcm(0, 5)
	expected := 0

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestLcmWithAllZeros(t *testing.T) {
	actual := helper.Lcm(0, 0)
	expected := 0

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

// TestLcmAvoidsMultiplyBeforeDivideOverflow uses values whose product
// overflows a 64-bit int (16000000016000000000, larger than
// math.MaxInt64) even though their LCM (4000000004000000000) does not.
// Computing (values[i]*lcm)/Gcd(...) would overflow and silently wrap
// to a wrong, negative result; dividing first avoids it.
func TestLcmAvoidsMultiplyBeforeDivideOverflow(t *testing.T) {
	actual := helper.Lcm(4000000000, 4000000004)
	expected := 4000000004000000000

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}
