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

func TestT3Simple(t *testing.T) {
	prices := make([]float64, 400)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}

	input := helper.SliceToChan(prices)
	t3 := trend.NewT3[float64]()
	result := t3.Compute(input)

	resultSlice := helper.ChanToSlice(result)

	expected := len(prices) - t3.IdlePeriod()
	if len(resultSlice) != expected {
		t.Fatalf("expected %d values, got %d", expected, len(resultSlice))
	}

	for i, v := range resultSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("T3 at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("T3: %d values, first: %.2f", len(resultSlice), resultSlice[0])
}

// referenceEma reproduces Ema's algorithm (SMA-seeded start, then the
// standard EMA recurrence) over a plain slice, to independently verify
// T3's values rather than just its output length.
func referenceEma(input []float64, period int) []float64 {
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += input[i]
	}

	prev := sum / float64(period)
	out := []float64{prev}

	multiplier := 2.0 / float64(period+1)
	for i := period; i < len(input); i++ {
		prev = (input[i]-prev)*multiplier + prev
		out = append(out, prev)
	}

	return out
}

// TestT3Values guards against a regression where ema3, ema4, and ema5 were
// each read by two consumers (the next EMA in the chain, and the final
// weighted sum) without being duplicated first, and combined with ema6
// without aligning for the extra EMA passes each is ahead by. Both silently
// produced wrong output rather than a build or panic, so this compares
// against an independently computed reference instead of just checking
// length or NaN/Inf, the way TestT3Simple already did.
func TestT3Values(t *testing.T) {
	prices := make([]float64, 400)
	for i := range prices {
		prices[i] = 100 + float64(i%7)
	}

	period := trend.DefaultT3Period
	factor := trend.DefaultT3VolumeFactor

	t3 := trend.NewT3[float64]()
	out := helper.ChanToSlice(t3.Compute(helper.SliceToChan(prices)))

	expectedLen := len(prices) - t3.IdlePeriod()
	if len(out) != expectedLen {
		t.Fatalf("expected %d values, got %d", expectedLen, len(out))
	}

	ema3 := referenceEma(referenceEma(referenceEma(prices, period), period), period)
	ema4 := referenceEma(ema3, period)
	ema5 := referenceEma(ema4, period)
	ema6 := referenceEma(ema5, period)

	a := factor
	a2, a3 := a*a, a*a*a
	c1, c2, c3, c4 := -a3, 3*a2+3*a3, -6*a2-3*a-3*a3, 1+3*a+3*a2+a3

	tail := func(s []float64) []float64 { return s[len(s)-expectedLen:] }
	e3, e4, e5, e6 := tail(ema3), tail(ema4), tail(ema5), tail(ema6)

	for i := 0; i < expectedLen; i++ {
		want := c1*e6[i] + c2*e5[i] + c3*e4[i] + c4*e3[i]
		if math.Abs(want-out[i]) > 1e-9 {
			t.Fatalf("value %d: expected %v, got %v", i, want, out[i])
		}
	}
}

// TestT3FlatInput guards against a regression where the weighted-sum
// coefficients didn't sum to 1 (c1+c2+c3+c4 = 3a(a-1) instead of 1 for the
// old, wrong coefficients), which meant T3 didn't track price at all -- a
// flat 100 series came out as -63. A moving average's weights must sum to
// 1, so a constant input should reproduce that same constant.
func TestT3FlatInput(t *testing.T) {
	const flat = 100.0

	prices := make([]float64, 100)
	for i := range prices {
		prices[i] = flat
	}

	t3 := trend.NewT3[float64]()
	out := helper.ChanToSlice(t3.Compute(helper.SliceToChan(prices)))

	if len(out) == 0 {
		t.Fatal("expected at least one value")
	}

	for i, v := range out {
		if math.Abs(v-flat) > 1e-9 {
			t.Fatalf("value %d: expected %v (flat input), got %v", i, flat, v)
		}
	}
}

func TestT3String(t *testing.T) {
	t3 := trend.NewT3[float64]()
	expected := "T3(5, 0.7)"
	actual := t3.String()
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestT3IdlePeriod(t *testing.T) {
	t3 := trend.NewT3[float64]()
	// 6 * (Period - 1) = 6 * 4 = 24
	expected := 24
	actual := t3.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
