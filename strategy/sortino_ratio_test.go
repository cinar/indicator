// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestSortinoRatioWithContext(t *testing.T) {
	// Equity curve 1, 1.5, 2.25, 3.375, 1.6875 corresponds to period returns of exactly
	// 0.5, 0.5, 0.5, -0.5 (each a clean power-of-two fraction, so no floating-point rounding
	// noise), giving a mean of 0.25. Only the last return falls below the zero minimum
	// acceptable return, giving a downside deviation of sqrt(0.5^2/4) = 0.25.
	outcomes := helper.SliceToChan([]float64{0, 0.5, 1.25, 2.375, 0.6875})

	actual := strategy.SortinoRatioWithContext(context.Background(), outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)
	expected := 1.0 * math.Sqrt(float64(strategy.DefaultSortinoRatioPeriodsPerYear))

	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSortinoRatioWithContextNoDownside(t *testing.T) {
	// Equity doubling every period (1, 2, 4, 8) gives an exactly representable return series of
	// 1.0, 1.0, 1.0, none of which fall below the zero minimum acceptable return, so the downside
	// deviation is zero.
	outcomes := helper.SliceToChan([]float64{0, 1, 3, 7})

	actual := strategy.SortinoRatioWithContext(context.Background(), outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSortinoRatioWithContextEmpty(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{})

	actual := strategy.SortinoRatioWithContext(context.Background(), outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSortinoRatioWithContextSingleOutcome(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{0})

	actual := strategy.SortinoRatioWithContext(context.Background(), outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSortinoRatioWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	outcomes := make(chan float64)

	cancel()

	done := make(chan float64)

	go func() {
		done <- strategy.SortinoRatioWithContext(ctx, outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)
	}()

	select {
	case actual := <-done:
		if actual != 0 {
			t.Fatalf("actual %v expected 0", actual)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("timeout - SortinoRatioWithContext did not return after cancellation")
	}
}

func TestSortinoRatio(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{0, 0.5, 1.25, 2.375, 0.6875})

	actual := strategy.SortinoRatio(outcomes, strategy.DefaultSortinoRatioPeriodsPerYear)
	expected := 1.0 * math.Sqrt(float64(strategy.DefaultSortinoRatioPeriodsPerYear))

	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
