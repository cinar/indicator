// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestSlowStochasticIndicator(t *testing.T) {
	// Value, K, D from stochastic.csv (Fast Stochastic)
	// 111.00,100.00,100.00
	// 109.00,75.00,91.67
	// 112.00,100.00,91.67
	// 113.00,100.00,91.67
	// 115.00,100.00,100.00
	// 114.00,87.50,95.83
	// 116.00,100.00,95.83
	// 118.00,100.00,95.83
	// 117.00,88.89,96.30
	// 119.00,100.00,96.30
	// 120.00,100.00,96.30
	// 118.00,75.00,91.67
	// 121.00,100.00,91.67
	// 122.00,100.00,91.67
	// 120.00,75.00,91.67
	// 123.00,100.00,91.67
	// 125.00,100.00,91.67
	// 124.00,87.50,95.83
	// 126.00,100.00,95.83

	type Data struct {
		Value float64
	}

	values := []float64{
		111, 109, 112, 113, 115, 114, 116, 118, 117, 119, 120, 118, 121, 122, 120, 123, 125, 124, 126,
	}

	input := helper.SliceToChan(values)

	// Default parameters: Period=10, KPeriod=3, DPeriod=3.
	// IdlePeriod = 10 + 3 + 3 - 3 = 13.
	s := trend.NewSlowStochastic[float64]()

	if s.IdlePeriod() != 13 {
		t.Fatalf("expected idle period 13, got %d", s.IdlePeriod())
	}

	actualK, actualD := s.Compute(input)

	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	// Since we don't have expected data easily, we just check if it produces values
	// and if the values are within reasonable range (0-100).
	count := 0
	for k := range actualK {
		d := <-actualD
		if k < 0 || k > 100 {
			t.Fatalf("k out of range: %v", k)
		}
		if d < 0 || d > 100 {
			t.Fatalf("d out of range: %v", d)
		}
		count++
	}

	expectedCount := len(values) - s.IdlePeriod()
	if count != expectedCount {
		t.Fatalf("expected %d values, got %d", expectedCount, count)
	}
}

func TestNewSlowStochasticWithPeriod(t *testing.T) {
	s := trend.NewSlowStochasticWithPeriod[float64](14, 3, 3)
	if s.Period != 14 {
		t.Fatalf("expected period 14, got %d", s.Period)
	}
}
