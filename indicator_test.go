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
