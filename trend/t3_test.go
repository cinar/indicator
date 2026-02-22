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
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}

	input := helper.SliceToChan(prices)
	t3 := trend.NewT3[float64]()
	result := t3.Compute(input)

	resultSlice := helper.ChanToSlice(result)

	if len(resultSlice) == 0 {
		t.Fatal("T3 produced no output")
	}

	for i, v := range resultSlice {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("T3 at index %d is NaN or Inf: %v", i, v)
		}
	}

	t.Logf("T3: %d values, first: %.2f", len(resultSlice), resultSlice[0])
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
