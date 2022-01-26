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
