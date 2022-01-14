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
