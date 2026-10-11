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

func TestSharpeRatioWithContext(t *testing.T) {
	// Equity curve 1, 1.1, 1.43 corresponds to period returns of exactly
	// 0.1 and 0.3, giving a mean of 0.2 and a population standard
	// deviation of 0.1.
	outcomes := helper.SliceToChan([]float64{0, 0.1, 0.43})

	actual := strategy.SharpeRatioWithContext(context.Background(), outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)
	expected := 2.0 * math.Sqrt(float64(strategy.DefaultSharpeRatioPeriodsPerYear))

	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSharpeRatioWithContextZeroVariance(t *testing.T) {
	// Equity doubling every period (1, 2, 4, 8) gives an exactly representable, zero-variance
	// return series of 1.0, 1.0, 1.0, avoiding floating-point rounding noise.
	outcomes := helper.SliceToChan([]float64{0, 1, 3, 7})

	actual := strategy.SharpeRatioWithContext(context.Background(), outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSharpeRatioWithContextEmpty(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{})

	actual := strategy.SharpeRatioWithContext(context.Background(), outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSharpeRatioWithContextSingleOutcome(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{0})

	actual := strategy.SharpeRatioWithContext(context.Background(), outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)

	if actual != 0 {
		t.Fatalf("actual %v expected 0", actual)
	}
}

func TestSharpeRatioWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	outcomes := make(chan float64)

	cancel()

	done := make(chan float64)

	go func() {
		done <- strategy.SharpeRatioWithContext(ctx, outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)
	}()

	select {
	case actual := <-done:
		if actual != 0 {
			t.Fatalf("actual %v expected 0", actual)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("timeout - SharpeRatioWithContext did not return after cancellation")
	}
}

func TestSharpeRatio(t *testing.T) {
	outcomes := helper.SliceToChan([]float64{0, 0.1, 0.43})

	actual := strategy.SharpeRatio(outcomes, strategy.DefaultSharpeRatioPeriodsPerYear)
	expected := 2.0 * math.Sqrt(float64(strategy.DefaultSharpeRatioPeriodsPerYear))

	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
