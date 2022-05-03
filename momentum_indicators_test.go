// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestChaikinOscillator(t *testing.T) {
	high := []float64{10, 11, 12, 13, 14, 15, 16, 17}
	low := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	closing := []float64{5, 6, 7, 8, 9, 10, 11, 12}
	volume := []int64{100, 200, 300, 400, 500, 600, 700, 800}
	expected := []float64{0, -7.41, -18.52, -31.69, -46.09, -61.27, -76.95, -92.97}

	actual, _ := ChaikinOscillator(2, 5, low, high, closing, volume)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestStochasticOscillator(t *testing.T) {
	high := []float64{
		127.01, 127.62, 126.59, 127.35, 128.17,
		128.43, 127.37, 126.42, 126.90, 126.85,
		125.65, 125.72, 127.16, 127.72, 127.69,
		128.22, 128.27, 128.09, 128.27, 127.74,
	}
	low := []float64{
		125.36, 126.16, 124.93, 126.09, 126.82,
		126.48, 126.03, 124.83, 126.39, 125.72,
		124.56, 124.57, 125.07, 126.86, 126.63,
		126.80, 126.71, 126.80, 126.13, 125.92,
	}
	closing := []float64{
		126.00, 126.60, 127.10, 127.20, 128.10,
		128.20, 126.30, 126.00, 126.60, 127.00,
		127.50, 128.00, 128.10, 127.29, 127.18,
		128.01, 127.11, 127.73, 127.06, 127.33,
	}
	expectedK := []float64{
		38.79, 54.87, 80.67, 84.39, 97.84,
		93.43, 39.14, 32.5, 49.17, 60.28,
		75.97, 88.89, 91.47, 70.54, 67.7,
		89.15, 65.89, 81.91, 64.60, 74.66,
	}

	actualK, _ := StochasticOscillator(high, low, closing)
	testEquals(t, roundDigitsAll(actualK, 2), expectedK)
}
