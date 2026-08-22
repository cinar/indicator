// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"math"
	"math/rand"
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
// confirm the compensated-summation implementation is meaningfully more
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

// TestMovingSumRecoversFromNaN guards against a regression where a single
// NaN entering the running sum poisoned every output for the rest of the
// series, even long after that NaN left the window: subtracting a value
// back out (`sum - NaN`) doesn't undo NaN contamination, since the running
// sum itself had already become NaN. The output should only be NaN while
// the NaN value is actually within the window, and recover to a correct,
// finite sum as soon as it slides out.
func TestMovingSumRecoversFromNaN(t *testing.T) {
	const period = 4

	values := []float64{1, 2, 3, math.NaN(), 4, 5, 6, 7, 8, 9}

	sum := trend.NewMovingSum[float64]()
	sum.Period = period

	actual := helper.ChanToSlice(sum.Compute(helper.SliceToChan(values)))

	expectedLen := len(values) - (period - 1)
	if len(actual) != expectedLen {
		t.Fatalf("expected %d values, got %d", expectedLen, len(actual))
	}

	for i, v := range actual {
		// Window i covers values[i : i+period]; NaN is at index 3.
		wantNaN := i <= 3 && i+period > 3

		if wantNaN {
			if !math.IsNaN(v) {
				t.Fatalf("value %d: expected NaN (window still contains the NaN input), got %v", i, v)
			}

			continue
		}

		var want float64
		for j := i; j < i+period; j++ {
			want += values[j]
		}

		if math.Abs(v-want) > 1e-9 {
			t.Fatalf("value %d: expected %v once the NaN left the window, got %v", i, want, v)
		}
	}
}

// TestMovingSumFloatDriftNearZeroSum exercises the regime where the window
// sum is near zero while individual terms are much larger in magnitude —
// the shape of a Cmf or Mfi moving sum of signed money flow in a quiet
// market. This is the case where classic Kahan summation is known to be
// weaker than plain summation; Neumaier's variant, used here, handles it
// correctly.
func TestMovingSumFloatDriftNearZeroSum(t *testing.T) {
	const (
		n      = 200_000
		period = 20
	)

	rnd := rand.New(rand.NewSource(1))

	values := make([]float64, n)
	for i := range values {
		values[i] = (rnd.Float64()*2 - 1) * 1e6
	}

	sum := trend.NewMovingSum[float64]()
	sum.Period = period

	actual := helper.ChanToSlice(sum.Compute(helper.SliceToChan(values)))

	var (
		naive        float64
		naiveMaxErr  float64
		actualMaxErr float64
	)

	for i, c := range values {
		var b float64
		if i-period >= 0 {
			b = values[i-period]
		}

		naive = naive + c - b

		if i >= period-1 {
			// Reference: the window summed from scratch, independent of
			// either running-sum implementation's accumulated state.
			var reference float64
			for j := i - period + 1; j <= i; j++ {
				reference += values[j]
			}

			actualErr := math.Abs(actual[i-period+1] - reference)
			if actualErr > actualMaxErr {
				actualMaxErr = actualErr
			}

			naiveErr := math.Abs(naive - reference)
			if naiveErr > naiveMaxErr {
				naiveMaxErr = naiveErr
			}
		}
	}

	if actualMaxErr >= naiveMaxErr {
		t.Fatalf("MovingSum (%v) is not more accurate than the naive running sum (%v) in the near-zero-sum regime", actualMaxErr, naiveMaxErr)
	}
}
