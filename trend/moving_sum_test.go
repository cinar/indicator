// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMovingSum(t *testing.T) {
	input := helper.SliceToChan([]int{-10, 20, -4, -5, 1, 5, 8, 10, -20, 4})
	expected := helper.SliceToChan([]int{1, 12, -3, 9, 24, 3, 2})

	sum := trend.NewMovingSum[int]()
	sum.Period = 4

	actual := sum.Compute(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// TestMovingSumFloatDrift checks that the running sum stays close to a
// freshly recomputed window sum over a long float32 series, rather than
// letting per-step rounding error compound as the series grows. It also
// computes the naive `sum = sum + c - b` running sum for the same series to
// confirm the Kahan-compensated implementation is meaningfully more
// accurate, not merely within some arbitrary tolerance.
func TestMovingSumFloatDrift(t *testing.T) {
	const (
		n      = 200_000
		period = 20
	)

	values := make([]float32, n)
	for i := range values {
		trendValue := 100.0 + float64(i)*0.01
		noise := math.Sin(float64(i)) * 5

		values[i] = float32(trendValue + noise)
	}

	sum := trend.NewMovingSum[float32]()
	sum.Period = period

	actual := helper.ChanToSlice(sum.Compute(helper.SliceToChan(values)))

	var (
		naive        float32
		reference    float64
		naiveMaxErr  float64
		actualMaxErr float64
	)

	for i, c := range values {
		var b float32
		if i-period >= 0 {
			b = values[i-period]
		}

		naive = naive + c - b

		var bRef float64
		if i-period >= 0 {
			bRef = float64(values[i-period])
		}

		reference = reference + float64(c) - bRef

		if i >= period-1 {
			actualErr := math.Abs(float64(actual[i-period+1]) - reference)
			if actualErr > actualMaxErr {
				actualMaxErr = actualErr
			}

			naiveErr := math.Abs(float64(naive) - reference)
			if naiveErr > naiveMaxErr {
				naiveMaxErr = naiveErr
			}
		}
	}

	const maxAllowedErr = 0.01

	if actualMaxErr > maxAllowedErr {
		t.Fatalf("MovingSum drifted too far from the reference sum: %v (max allowed %v)", actualMaxErr, maxAllowedErr)
	}

	if actualMaxErr >= naiveMaxErr {
		t.Fatalf("MovingSum (%v) is not more accurate than the naive running sum (%v)", actualMaxErr, naiveMaxErr)
	}
}
