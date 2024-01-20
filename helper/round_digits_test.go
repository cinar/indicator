// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRoundDigits(t *testing.T) {
	input := []float64{10.1234, 5.678, 6.78, 8.91011}
	expected := []float64{10.12, 5.68, 6.78, 8.91}

	actual := helper.ChanToSlice(helper.RoundDigits(helper.SliceToChan(input), 2))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
