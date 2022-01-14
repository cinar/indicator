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
