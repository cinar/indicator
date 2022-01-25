// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import "testing"

func TestMoneyFlowIndex(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{100, 100, 57.01, 65.85, 61.54}
	period := 2

	actual := roundDigitsAll(MoneyFlowIndex(period, high, low, closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestMoneyFlowIndex2(t *testing.T) {
	high := []float64{2390.9, 2386.3, 2395.33, 2399.0, 2402.46, 2401.15, 2421.98, 2430.31, 2426.33, 2434.93, 2470.83, 2483.36, 2467.19, 2450.72}
	low := []float64{2373.15, 2370.0, 2380.77, 2384.28, 2387.46, 2385.02, 2383.18, 2408.39, 2410.59, 2420.89, 2428.92, 2456.77, 2437.65, 2440.87}
	closing := []float64{2373.39, 2382.47, 2394.4, 2387.51, 2395.64, 2389.47, 2410.24, 2425.37, 2422.33, 2430.29, 2465.76, 2466.27, 2440.07, 2445.85}
	volume := []int64{1621, 1387, 1444, 1298, 1629, 1598, 2311, 2934, 2128, 1823, 5078, 6693, 3960, 1927}
	expected := []float64{100, 53.884998, 68.888204, 53.300054, 63.645476, 52.296425, 62.119198, 70.011701, 60.826077, 54.659774, 64.728457, 72.749002, 64.185176, 60.710993}
	period := 14

	actual := roundDigitsAll(MoneyFlowIndex(period, high, low, closing, volume), 6)
	testEquals(t, actual, expected)
}

func TestForceIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{900, 220, -320, 360, -180}
	period := 1

	actual := roundDigitsAll(ForceIndex(period, closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestDefaultForceIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{900, 802.86, 642.45, 602.1, 490.37}

	actual := roundDigitsAll(DefaultForceIndex(closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestDefaultEaseOfMovement(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{32000000, 16000000, 13791666.67, 11385416.67, 8219444.44}

	actual := roundDigitsAll(DefaultEaseOfMovement(high, low, volume), 2)
	testEquals(t, actual, expected)
}

func TestVolumePriceTrend(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{0, 24.44, -4.65, 46.78, 28.78}

	actual := roundDigitsAll(VolumePriceTrend(closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestVolumeWeightedAveragePrice(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	period := 2
	expected := []float64{9, 10.05, 9.32, 8.8, 9.14}

	actual := roundDigitsAll(VolumeWeightedAveragePrice(period, closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestDefaultVolumeWeightedAveragePrice(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{9, 10.05, 9.21, 9.44, 9.18}

	actual := roundDigitsAll(DefaultVolumeWeightedAveragePrice(closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestNegativeVolumeIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{1000, 1000, 636.36, 636.36, 509.09}

	actual := roundDigitsAll(NegativeVolumeIndex(closing, volume), 2)
	testEquals(t, actual, expected)
}

func TestChaikinMoneyFlow(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	volume := []int64{100, 110, 80, 120, 90}
	closing := []float64{9, 11, 7, 10, 8}
	expected := []float64{0.5, 1.81, 0.67, -0.41, -0.87}

	actual := roundDigitsAll(ChaikinMoneyFlow(high, low, closing, volume), 2)
	testEquals(t, actual, expected)
}
