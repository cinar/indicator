// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestStcSlowStochastic(t *testing.T) {
	resultChan := make(chan float64, 10)
	go func() {
		values := []float64{
			100, 102, 101, 103, 105, 104, 106, 108, 107, 109,
			110, 111, 109, 112, 113, 115, 114, 116, 118, 117,
			119, 120, 118, 121, 122, 120, 123, 125, 124, 126,
		}

		input := helper.SliceToChan(values)
		s := trend.NewSlowStochastic[float64]()
		actualK, actualD := s.Compute(input)

		for {
			select {
			case v, ok := <-actualK:
				if !ok {
					return
				}
				resultChan <- v
			case v, ok := <-actualD:
				if !ok {
					return
				}
				resultChan <- v
			case <-time.After(2 * time.Second):
				return
			}
		}
	}()

	select {
	case v := <-resultChan:
		t.Logf("got value: %v", v)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout - no values produced")
	}
}

// TestStcFull checks that STC's output length matches IdlePeriod() once the
// Stochastic pass is applied twice (once to MACD, again to its %D), which
// needs more warm-up than the 19-row testdata/stochastic.csv fixture
// (shared with TestStcSlowStochastic/TestStochastic) has, so this generates
// its own longer synthetic series instead.
func TestStcFull(t *testing.T) {
	values := make([]float64, 200)
	for i := range values {
		values[i] = 100 + float64(i%7)
	}

	stc := trend.NewStcWithPeriod[float64](5, 10, 5, 3)
	result := stc.Compute(helper.SliceToChan(values))

	slice := helper.ChanToSlice(result)

	expected := len(values) - stc.IdlePeriod()
	if len(slice) != expected {
		t.Fatalf("expected %d values, got %d", expected, len(slice))
	}
}

// TestStcDeadlock guards against a regression where ComputeWithContext's
// output channel never closed: %K was consumed by two independent Subtract
// stages sharing the same channel, instead of a duplicated one each, which
// deadlocked the surrounding duplicate/buffer chain once the series was
// long enough (the 20-row fixture in TestStcFull is too short to trigger
// it -- this needs several hundred points).
func TestStcDeadlock(t *testing.T) {
	closings := make([]float64, 400)
	for i := range closings {
		closings[i] = 100 + float64(i%7)
	}

	stc := trend.NewStc[float64]()

	done := make(chan []float64, 1)
	go func() {
		done <- helper.ChanToSlice(stc.Compute(helper.SliceToChan(closings)))
	}()

	select {
	case slice := <-done:
		expected := len(closings) - stc.IdlePeriod()
		if len(slice) != expected {
			t.Fatalf("expected %d values, got %d", expected, len(slice))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Stc.ComputeWithContext deadlocked")
	}
}
