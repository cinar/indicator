// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSafeDivide(t *testing.T) {
	input := 4.0
	divider := 2.0
	fallback := 0.5
	expected := 2.0

	actual := helper.SafeDivide(input, divider, fallback)

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSafeDivideByZero(t *testing.T) {
	input := 4.0
	divider := 0.0
	fallback := 0.5
	expected := fallback

	actual := helper.SafeDivide(input, divider, fallback)

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
