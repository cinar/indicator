// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestGcdWithTwoValues(t *testing.T) {
	actual := helper.Gcd(1220, 512)
	expected := 4

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestGcdWithThreeValues(t *testing.T) {
	actual := helper.Gcd(1, 2, 5)
	expected := 1

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestGcdWithFourValues(t *testing.T) {
	actual := helper.Gcd(2, 4, 6, 12)
	expected := 2

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}
