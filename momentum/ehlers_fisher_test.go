// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestEhlersFisherSimple(t *testing.T) {
	highs := make([]float64, 400)
	lows := make([]float64, 400)

	for i := range highs {
		highs[i] = 101 + float64(i)
		lows[i] = 99 + float64(i)
	}

	ehlersFisher := momentum.NewEhlersFisher[float64]()
	result := ehlersFisher.Compute(helper.SliceToChan(highs), helper.SliceToChan(lows))

	resultSlice := helper.ChanToSlice(result)

	expected := len(highs) - ehlersFisher.IdlePeriod()
	if len(resultSlice) != expected {
		t.Fatalf("expected %d values, got %d", expected, len(resultSlice))
	}

	for i, v := range resultSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("EhlersFisher at index %d is NaN or Inf: %v", i, v)
		}
	}
}

// referenceEhlersFisher computes the canonical, recursive Fisher Transform
// with a plain sliding-window min/max over slices, independently of
// EhlersFisher's channel implementation, following the formula in Ehlers'
// "Using the Fisher Transform" article.
func referenceEhlersFisher(highs, lows []float64, period int) []float64 {
	prices := make([]float64, len(highs))
	for i := range highs {
		prices[i] = (highs[i] + lows[i]) / 2
	}

	var out []float64

	value1Prev := 0.0
	fisherPrev := 0.0

	for i := period - 1; i < len(prices); i++ {
		window := prices[i-period+1 : i+1]

		lo, hi := window[0], window[0]
		for _, v := range window {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}

		normalized := (prices[i] - lo) / (hi - lo)

		value1 := 0.33*2*(normalized-0.5) + 0.67*value1Prev
		if value1 > 0.99 {
			value1 = 0.999
		}
		if value1 < -0.99 {
			value1 = -0.999
		}

		fisher := 0.5*math.Log((1+value1)/(1-value1)) + 0.5*fisherPrev

		out = append(out, fisher)

		value1Prev = value1
		fisherPrev = fisher
	}

	return out
}

// referenceEhlersFisherWithoutRecursion computes the same formula but with
// the previous Value1/Fisher terms omitted (always treated as 0), as if the
// recursive smoothing had been accidentally left out. It exists solely so
// TestEhlersFisherRecursionIsApplied can confirm EhlersFisher's real output
// is not accidentally equivalent to this non-recursive shortcut.
func referenceEhlersFisherWithoutRecursion(highs, lows []float64, period int) []float64 {
	prices := make([]float64, len(highs))
	for i := range highs {
		prices[i] = (highs[i] + lows[i]) / 2
	}

	var out []float64

	for i := period - 1; i < len(prices); i++ {
		window := prices[i-period+1 : i+1]

		lo, hi := window[0], window[0]
		for _, v := range window {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}

		normalized := (prices[i] - lo) / (hi - lo)

		value1 := 0.33 * 2 * (normalized - 0.5)
		if value1 > 0.99 {
			value1 = 0.999
		}
		if value1 < -0.99 {
			value1 = -0.999
		}

		fisher := 0.5 * math.Log((1+value1)/(1-value1))

		out = append(out, fisher)
	}

	return out
}

// TestEhlersFisherHandComputedValues cross-checks EhlersFisher's output,
// row by row, against TA-Lib FISHER(high, low, 5) for the same input.
func TestEhlersFisherHandComputedValues(t *testing.T) {
	highs := []float64{10, 12, 11, 13, 15, 14, 16, 18, 17, 19, 21, 20, 22, 24, 23}
	lows := []float64{8, 9, 10, 11, 12, 11, 13, 15, 14, 16, 17, 16, 18, 20, 19}
	period := 5

	// Values produced by TA-Lib 0.8.3 FISHER for this exact input series
	// and period.
	expected := []float64{
		0.342828254415,
		0.515477314672,
		0.878757489300,
		1.306151604525,
		1.400850729870,
		1.684116084891,
		2.048541828445,
		1.894115311568,
		2.045302533791,
		2.339153146801,
		2.131560751273,
	}

	if len(expected) < 8 {
		t.Fatalf("test fixture too small to be a meaningful cross-check")
	}

	ehlersFisher := momentum.NewEhlersFisherWithPeriod[float64](period)
	actual := helper.ChanToSlice(ehlersFisher.Compute(helper.SliceToChan(highs), helper.SliceToChan(lows)))

	if len(actual) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if math.Abs(actual[i]-expected[i]) > 1e-9 {
			t.Fatalf("value %d: expected %v, got %v", i, expected[i], actual[i])
		}
	}

	// Sanity check the hand-computed reference implementation above
	// against the same independently-verified values.
	ref := referenceEhlersFisher(highs, lows, period)
	for i := range expected {
		if math.Abs(ref[i]-expected[i]) > 1e-9 {
			t.Fatalf("reference value %d: expected %v, got %v", i, expected[i], ref[i])
		}
	}
}

// TestEhlersFisherRecursionIsApplied confirms that EhlersFisher actually
// carries Value1 and Fisher across iterations, rather than accidentally
// recomputing each output independently (as referenceEhlersFisherWithoutRecursion
// does). If the recursive smoothing were dropped, EhlersFisher's output
// would match referenceEhlersFisherWithoutRecursion and this test would fail.
func TestEhlersFisherRecursionIsApplied(t *testing.T) {
	highs := []float64{10, 12, 11, 13, 15, 14, 16, 18, 17, 19, 21, 20, 22, 24, 23}
	lows := []float64{8, 9, 10, 11, 12, 11, 13, 15, 14, 16, 17, 16, 18, 20, 19}
	period := 5

	ehlersFisher := momentum.NewEhlersFisherWithPeriod[float64](period)
	actual := helper.ChanToSlice(ehlersFisher.Compute(helper.SliceToChan(highs), helper.SliceToChan(lows)))

	withoutRecursion := referenceEhlersFisherWithoutRecursion(highs, lows, period)

	if len(actual) != len(withoutRecursion) {
		t.Fatalf("expected %d values, got %d", len(withoutRecursion), len(actual))
	}

	// Index 0 coincides by construction: both formulas start from the same
	// zero-seeded Value1[previous]/Fisher[previous], so they only diverge
	// once a real previous value has been fed back in.
	for i := 1; i < len(actual); i++ {
		if math.Abs(actual[i]-withoutRecursion[i]) < 1e-6 {
			t.Fatalf("value %d: actual %v unexpectedly matches non-recursive %v; recursion may be missing", i, actual[i], withoutRecursion[i])
		}
	}
}

func TestEhlersFisherString(t *testing.T) {
	ehlersFisher := momentum.NewEhlersFisher[float64]()
	expected := "EhlersFisher(10)"
	actual := ehlersFisher.String()
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestEhlersFisherIdlePeriod(t *testing.T) {
	ehlersFisher := momentum.NewEhlersFisher[float64]()
	expected := 9
	actual := ehlersFisher.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestEhlersFisher(t *testing.T) {
	type Data struct {
		High   float64
		Low    float64
		Fisher float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/ehlers_fisher.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })

	ehlersFisher := momentum.NewEhlersFisher[float64]()
	actual := ehlersFisher.Compute(highs, lows)

	actual = helper.RoundDigits(actual, 6)

	inputs[2] = helper.Skip(inputs[2], ehlersFisher.IdlePeriod())

	err = helper.CheckEquals(
		helper.Map(actual, func(v float64) float64 { return v }),
		helper.Map(inputs[2], func(d *Data) float64 { return d.Fisher }),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TestEhlersFisherClamp checks the Value1 clamp, which sets Value1 to
// ±0.999 once it passes ±0.99, against TA-Lib FISHER(high, low, 5) on a
// series that rises and then falls steadily enough to trigger it.
func TestEhlersFisherClamp(t *testing.T) {
	var highs, lows []float64
	for i := 0; i < 40; i++ {
		mid := 100 + float64(i)
		if i >= 20 {
			mid = 100 + float64(38-i)
		}

		highs = append(highs, mid+1)
		lows = append(lows, mid-1)
	}

	// Values produced by TA-Lib 0.8.3 FISHER for this exact input series
	// and period. The jumps after 4.686 and after -4.507 are the clamp.
	expected := []float64{
		0.342828254415, 0.791373872129, 1.261492949351, 1.725174983451,
		2.175408580680, 2.612568601427, 3.039143272643, 3.457799327674,
		3.870773984113, 4.279770426648, 4.686022806084, 6.143212570292,
		6.871807452396, 7.236104893448, 7.418253613974, 7.509327974237,
		4.798325884638, 2.593732724946, 1.092822252263, 0.042939059260,
		-0.739155652702, -1.365244343253, -1.900559549376, -2.382692314323,
		-2.833314002311, -3.265009971831, -3.685172921140, -4.098205747151,
		-4.506770165842, -6.053586250171, -6.826994292336, -7.213698313418,
		-7.407050323959, -7.503726329230, -7.552064331865, -7.576233333183,
	}

	ehlersFisher := momentum.NewEhlersFisherWithPeriod[float64](5)
	actual := helper.ChanToSlice(ehlersFisher.Compute(helper.SliceToChan(highs), helper.SliceToChan(lows)))

	if len(actual) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if math.Abs(actual[i]-expected[i]) > 1e-9 {
			t.Fatalf("value %d: expected %v, got %v", i, expected[i], actual[i])
		}
	}
}
