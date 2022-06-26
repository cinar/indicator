// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestBollingerBands(t *testing.T) {
	closing := []float64{
		2, 4, 6, 8, 12, 14, 16, 18, 20,
		2, 4, 6, 8, 12, 14, 16, 18, 20,
		2, 4, 6, 8, 12, 14, 16, 18, 20,
		2, 4, 6, 8, 12, 14, 16, 18, 20,
	}

	expectedMiddleBand := []float64{
		2, 3, 4, 5, 6.4, 7.67, 8.86, 10, 11.11,
		10.2, 9.64, 9.33, 9.23, 9.43, 9.73, 10.13, 10.59, 11.11,
		10.63, 10.3, 10.5, 10.7, 11, 11.3, 11.5, 11.7, 11.9,
		11.1, 10.3, 10.5, 10.7, 11, 11.3, 11.5, 11.7, 11.9,
	}

	expectedUpperBand := []float64{
		2, 3, 4, 5, 6.4, 7.67, 8.86, 10, 11.11, 10.2,
		9.64, 9.33, 9.23, 9.43, 9.73, 10.13, 10.59, 11.11,
		10.63, 22.78, 22.56, 22.45, 22.56, 22.84, 23.22, 23.72, 24.32,
		23.9, 22.78, 22.56, 22.45, 22.56, 22.84, 23.22, 23.72, 24.32,
	}

	expectedLowerBand := []float64{
		2, 3, 4, 5, 6.4, 7.67, 8.86, 10, 11.11,
		10.2, 9.64, 9.33, 9.23, 9.43, 9.73, 10.13, 10.59, 11.11,
		10.63, -2.18, -1.56, -1.05, -0.56, -0.24, -0.22, -0.32, -0.52,
		-1.7, -2.18, -1.56, -1.05, -0.56, -0.24, -0.22, -0.32, -0.52,
	}

	actualMiddleBand, actualUpperBand, actualLowerBand := BollingerBands(closing)
	testEquals(t, roundDigitsAll(actualMiddleBand, 2), expectedMiddleBand)
	testEquals(t, roundDigitsAll(actualUpperBand, 2), expectedUpperBand)
	testEquals(t, roundDigitsAll(actualLowerBand, 2), expectedLowerBand)
}

func TestStd(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	expected := []float64{0, 1, 1, 1, 2, 1, 1, 1, 1}
	period := 2

	actual := Std(period, values)
	testEquals(t, roundDigitsAll(actual, 3), expected)
}

func TestStdFromSma(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	expected := []float64{0, 1, 1, 1, 2, 1, 1, 1, 1}
	period := 2

	actual := StdFromSma(period, values, Sma(period, values))
	testEquals(t, roundDigitsAll(actual, 3), expected)
}

func TestUlcerIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8, 7, 7, 8, 10, 9, 5, 4, 6, 7}
	expected := []float64{0, 0, 20.99, 18.74, 20.73, 24.05, 26.17, 26.31,
		24.99, 24.39, 28.49, 32.88, 34.02, 34.19}

	actual := DefaultUlcerIndex(closing)
	testEquals(t, roundDigitsAll(actual, 2), expected)
}

func TestDonchianChannel(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	period := 4
	expectedUpperChannel := []float64{9, 11, 11, 11, 11}
	expectedMiddleChannel := []float64{9, 10, 9, 9, 9}
	expectedLowerChannel := []float64{9, 9, 7, 7, 7}

	actualUpperChannel, actualMiddleChannel, actualLowerChannel := DonchianChannel(period, closing)
	testEquals(t, roundDigitsAll(actualUpperChannel, 2), expectedUpperChannel)
	testEquals(t, roundDigitsAll(actualMiddleChannel, 2), expectedMiddleChannel)
	testEquals(t, roundDigitsAll(actualLowerChannel, 2), expectedLowerChannel)
}

func TestKeltnerChannel(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	closing := []float64{9, 11, 7, 10, 8}
	expectedUpperBand := []float64{17, 17.19, 17.65, 17.58, 17.38}
	expectedMiddleLine := []float64{9, 9.19, 8.98, 9.08, 8.98}
	expectedLowerBand := []float64{1, 1.19, 0.32, 0.58, 0.58}

	actualUpperBand, actualMiddleLine, actualLowerBand := DefaultKeltnerChannel(high, low, closing)
	testEquals(t, roundDigitsAll(actualUpperBand, 2), expectedUpperBand)
	testEquals(t, roundDigitsAll(actualMiddleLine, 2), expectedMiddleLine)
	testEquals(t, roundDigitsAll(actualLowerBand, 2), expectedLowerBand)
}
