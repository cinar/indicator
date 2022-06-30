// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestAbsolutePriceOscillator(t *testing.T) {
	values := []float64{1, 2, 1, 5, 8, 10, 4, 6, 5, 2}
	expected := []float64{0, 0.33, 0, 1.26, 2.26, 2.65, 0.14, 0.22, -0.14, -1.19}

	actual := AbsolutePriceOscillator(2, 5, values)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestBalanceOfPower(t *testing.T) {
	opening := []float64{10, 20, 15, 50}
	high := []float64{40, 25, 20, 60}
	low := []float64{4, 10, 5, 6}
	closing := []float64{20, 15, 50, 55}
	expected := []float64{0.28, -0.33, 2.33, 0.09}

	actual := BalanceOfPower(opening, high, low, closing)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestChandeForecastOscillator(t *testing.T) {
	closing := []float64{1, 5, 12, 20}
	expected := []float64{110, -26, -5.8333, 4.5}

	actual := ChandeForecastOscillator(closing)
	testEquals(t, roundDigitsAll(actual, 4), expected)
}

func TestCommunityChannelIndex(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	closing := []float64{9, 11, 7, 10, 8}
	expected := []float64{0, 133.33, 114.29, 200, 26.32}

	actual := DefaultCommunityChannelIndex(high, low, closing)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestEma(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	expected := []float64{2, 3.333, 5.111, 7.037, 10.346, 12.782, 14.927, 16.976, 18.992}
	period := 2

	actual := Ema(period, values)
	testEquals(t, roundDigitsAll(actual, 3), expected)
}

func TestMassIndex(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	expected := []float64{1, 1.92, 2.83, 3.69, 4.52}

	actual := MassIndex(high, low)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestMax(t *testing.T) {
	values := []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	expected := []float64{10, 10, 10, 10, 9, 8, 7, 6, 5, 4}

	actual := Max(4, values)
	testEquals(t, actual, expected)
}

func TestMin(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []float64{1, 1, 1, 1, 2, 3, 4, 5, 6, 7}

	actual := Min(4, values)
	testEquals(t, actual, expected)
}

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

	closing := []float64{
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

	psar, trend := ParabolicSar(high, low, closing)
	testEquals(t, roundDigitsAll(psar, 2), expectedPsar)

	for i := 0; i < len(expectedTrend); i++ {
		if trend[i] != expectedTrend[i] {
			t.Fatalf("at %d actual %d expected %d", i, trend[i], expectedTrend[i])
		}
	}
}

func TestQstick(t *testing.T) {
	opening := []float64{10, 20, 15, 50, 40, 41, 43, 80}
	closing := []float64{20, 15, 50, 55, 42, 30, 31, 70}
	expected := []float64{10, 2.5, 13.33, 11.25, 9.4, 5.2, 3.8, -5.2}

	actual := Qstick(5, opening, closing)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestKdj(t *testing.T) {
	low := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	high := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	closing := []float64{5, 10, 15, 20, 25, 30, 35, 40, 45, 50}
	expectedK := []float64{44.44, 45.91, 46.70, 48.12, 48.66, 48.95, 49.14, 49.26, 49.36, 49.26}
	expectedD := []float64{44.44, 45.18, 45.68, 46.91, 47.82, 48.58, 48.91, 49.12, 49.25, 49.30}
	expectedJ := []float64{44.44, 47.37, 48.72, 50.55, 50.32, 49.70, 49.58, 49.56, 49.57, 49.19}

	k, d, j := DefaultKdj(high, low, closing)

	testEquals(t, roundDigitsAll(k, 2), expectedK)
	testEquals(t, roundDigitsAll(d, 2), expectedD)
	testEquals(t, roundDigitsAll(j, 2), expectedJ)
}

func TestRma(t *testing.T) {
	values := []float64{
		0,
		0.00005,
		0.000017,
		0.000262,
		0.000107,
		0,
		0,
		0.000597,
		0,
		0,
		0.000059,
		0.000198,
		0.000073,
		0,
		0.000006,
		0,
		0.000077,
		0.000032,
		0.000112,
	}

	expected := []float64{
		0.00009735714286,
		0.00009083163265,
		0.00008434365889,
		0.00008381911183,
		0.0000801177467,
		0.00008239505051,
	}
	period := 14

	actual := Rma(period, values)
	testEquals(t, roundDigitsAll(actual[len(actual)-6:], 14), expected)
}

func TestSma(t *testing.T) {
	values := []float64{2, 4, 6, 8, 10}
	expected := []float64{2, 3, 5, 7, 9}
	period := 2

	actual := Sma(period, values)
	testEquals(t, actual, expected)
}

func TestSince(t *testing.T) {
	values := []float64{1, 2, 2, 3, 4, 4, 4, 4, 5, 6}
	expected := []int{0, 0, 1, 0, 0, 1, 2, 3, 0, 0}

	actual := Since(values)
	testEqualsInt(t, actual, expected)
}

func TestSum(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []float64{1, 3, 6, 10, 14, 18, 22, 26, 30, 34}
	period := 4

	actual := Sum(period, values)
	testEquals(t, actual, expected)
}

func TestTrix(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	period := 4
	expected := []float64{0, 0.06, 0.17, 0.26, 0.33, 0.33, 0.3, 0.25, 0.21}

	actual := Trix(period, values)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestVortex(t *testing.T) {
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

	closing := []float64{
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

	plusVi, minusVi := Vortex(high, low, closing)
	testEquals(t, roundDigitsAll(plusVi, 5), expectedPlusVi)
	testEquals(t, roundDigitsAll(minusVi, 5), expectedMinusVi)
}

func TestVwma(t *testing.T) {
	closing := []float64{20, 21, 21, 19, 16}
	volume := []int64{100, 50, 40, 50, 100}
	expected := []float64{20, 20.33, 20.47, 20.29, 17.84}
	period := 3

	actual := Vwma(period, closing, volume)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestDefaultVwma(t *testing.T) {
	closing := []float64{20, 21, 21, 19, 16}
	volume := []int64{100, 50, 40, 50, 100}
	expected := []float64{20, 20.33, 20.47, 20.17, 18.94}

	actual := DefaultVwma(closing, volume)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}
