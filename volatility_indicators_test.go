// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestStd(t *testing.T) {
	values := []float64{2, 4, 6, 8, 12, 14, 16, 18, 20}
	expected := []float64{0, 0.707, 1, 1, 1.581, 1.581, 1, 1, 1}
	period := 2

	actual := Std(period, values)
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
