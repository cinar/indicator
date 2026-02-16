// Copyright (c) 2021-2026 Onur Cinar.
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
		22.51, 22.7, 22.88, 23.01, 23.06, 23.08, 23.16, 23.11, 23.15, 23.05,
		22.92, 22.81, 22.74, 22.66, 22.63, 22.74, 22.97, 23.12, 23.27, 23.43,
		23.52, 23.34,
	})

	dema := trend.NewDema[float64]()
	actual := dema.Compute(input)

	actual = helper.RoundDigits(actual, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
