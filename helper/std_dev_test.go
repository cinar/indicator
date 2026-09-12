// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestStdDev(t *testing.T) {
	actual := helper.StdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	expected := 2.0

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestStdDevEmpty(t *testing.T) {
	actual := helper.StdDev([]float64{})
	expected := 0.0

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestStdDevConstant(t *testing.T) {
	actual := helper.StdDev([]float64{3, 3, 3})
	expected := 0.0

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
