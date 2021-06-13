package indicator

import (
	"math"
	"testing"
)

func TestSma(t *testing.T) {
	values := []float64{2, 4, 6, 8, 10}
	sma := []float64{2, 3, 5, 7, 9}
	period := 2

	result := Sma(period, values)
	if len(result) != len(sma) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := result[i]
		expected := sma[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestStd(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	std := []float64{0, 0.707, 1, 1, 1.581, 1.581, 1, 1, 1}
	period := 2

	result := Std(period, values)
	if len(result) != len(values) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := math.Round(result[i]*1000) / 1000
		expected := std[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestEma(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	ema := []float64{2, 3.333, 5.111, 7.037, 10.346, 12.782, 14.927, 16.976, 18.992}
	period := 2

	result := Ema(period, values)
	if len(result) != len(ema) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := math.Round(result[i]*1000) / 1000
		expected := ema[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestSince(t *testing.T) {
	values := []float64{1, 2, 2, 3, 4, 4, 4, 4, 5, 6}
	expected := []int{0, 0, 1, 0, 0, 1, 2, 3, 0, 0}

	actual := Since(values)
	// TODO: check size.

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %d expected %d", i, actual[i], expected[i])
		}
	}
}
