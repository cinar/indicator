package indicator

import (
	"math"
	"testing"
)

func TestParabolicSAR(t *testing.T) {

	high := []float64{
		3836.86,
		3766.57,
		3576.17,
		3513.55,
		3529.75,
		3756.17,
		3717.17,
		3572.62,
		3612.43,
	}

	low := []float64{
		3643.25,
		3542.73,
		3371.75,
		3334.02,
		3314.75,
		3558.21,
		3517.79,
		3447.90,
		3494.39,
	}

	close := []float64{
		3790.55,
		3546.20,
		3507.31,
		3340.81,
		3529.60,
		3717.41,
		3544.35,
		3478.14,
		3612.08,
	}

	expectedPsar := []float64{
		3836.86,
		3836.86,
		3836.86,
		3808.95,
		3770.96,
		3314.75,
		3314.75,
		3323.58,
		3332.23,
	}

	expectedTrend := []Trend{
		Falling,
		Falling,
		Falling,
		Falling,
		Falling,
		Rising,
		Rising,
		Rising,
		Rising,
	}

	psar, trend := ParabolicSar(high, low, close)
	if len(psar) != len(expectedPsar) || len(trend) != len(expectedTrend) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expectedPsar); i++ {
		currentPsar := math.Round(psar[i]*100) / 100
		if currentPsar != expectedPsar[i] {
			t.Fatalf("at %d actual %f expected %f", i, currentPsar, expectedPsar[i])
		}

		if trend[i] != expectedTrend[i] {
			t.Fatalf("at %d actual %d expected %d", i, trend[i], expectedTrend[i])
		}
	}
}

func TestVertex(t *testing.T) {
	high := []float64{
		1404.14,
		1405.95,
		1405.98,
		1405.87,
		1410.03,
	}

	low := []float64{
		1396.13,
		1398.80,
		1395.62,
		1397.32,
		1400.60,
	}

	close := []float64{
		1402.22,
		1402.80,
		1405.87,
		1404.11,
		1403.93,
	}

	expectedPlusVi := []float64{
		0.00000,
		1.37343,
		0.97087,
		1.04566,
		1.12595,
	}

	expectedMinusVi := []float64{
		0.00000,
		0.74685,
		0.89492,
		0.93361,
		0.83404,
	}

	plusVi, minusVi := Vortex(high, low, close)
	if len(plusVi) != len(expectedPlusVi) || len(minusVi) != len(expectedMinusVi) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(plusVi); i++ {
		actualPlusVi := math.Round(plusVi[i]*100000) / 100000
		if actualPlusVi != expectedPlusVi[i] {
			t.Fatalf("at %d actual %f expected %f", i, actualPlusVi, expectedPlusVi[i])
		}

		actualMinusVi := math.Round(minusVi[i]*100000) / 100000
		if actualMinusVi != expectedMinusVi[i] {
			t.Fatalf("at %d actual %f expected %f", i, actualMinusVi, expectedMinusVi[i])
		}
	}
}

func TestChandeForecastOscillator(t *testing.T) {
	closing := []float64{1, 5, 12, 20}
	expected := []float64{110, -26, -5.8333, 4.5}

	actual := ChandeForecastOscillator(closing)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		a := roundDigits(actual[i], 4)

		if a != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, a, expected[i])
		}
	}
}
