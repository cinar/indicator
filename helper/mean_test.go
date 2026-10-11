// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestMean(t *testing.T) {
	actual := helper.Mean([]float64{1, 2, 3, 4})
	expected := 2.5

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestMeanEmpty(t *testing.T) {
	actual := helper.Mean([]float64{})
	expected := 0.0

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
