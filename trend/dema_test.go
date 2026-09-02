// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestDema(t *testing.T) {
	input := helper.SliceToChan([]float64{
		22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29,
		22.15, 22.39, 22.38, 22.61, 23.36, 24.05, 23.75, 23.83, 23.95, 23.63,
		23.82, 23.87, 23.65, 23.19, 23.10, 23.33, 22.68, 23.10, 22.40, 22.17,
		22.15, 22.39, 22.38, 22.61, 23.36, 24.05, 23.75, 23.83, 23.95, 23.63,
		22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24, 22.29,
		23.82, 23.87, 23.65, 23.19, 23.10, 23.33, 22.68, 23.10, 22.40, 22.17,
	})

	expected := helper.SliceToChan([]float64{
		23.38, 23.45, 23.26, 23.08, 22.9, 22.77, 22.65, 22.54, 22.47, 22.45,
		22.39, 22.35, 22.6, 22.82, 22.97, 23.02, 23.04, 23.1, 23.04, 23.05,
		22.94, 22.81,
	})

	dema := trend.NewDema[float64]()
	actual := dema.Compute(input)

	actual = helper.RoundDigits(actual, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDemaString(t *testing.T) {
	expected := "DEMA(20,20)"
	actual := trend.NewDema[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
