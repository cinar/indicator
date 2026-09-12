// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Mean returns the arithmetic mean of the given values, or zero for an empty slice.
//
// Example:
//
//	n := helper.Mean([]float64{1, 2, 3, 4})
//	fmt.Println(n) // 2.5
func Mean[T Float](values []T) T {
	if len(values) == 0 {
		return 0
	}

	var sum T

	for _, v := range values {
		sum += v
	}

	return sum / T(len(values))
}
