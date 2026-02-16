// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestDuplicate(t *testing.T) {
	expecteds := []float64{-10, 20, -4, -5}

	outputs := helper.Duplicate[float64](helper.SliceToChan(expecteds), 4)

	for i, expected := range expecteds {
		for _, output := range outputs {
			actual := <-output
			if actual != expected {
				t.Fatalf("index %d actual %v expected %v", i, actual, expected)
			}
		}
	}
}
