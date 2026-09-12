// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "math"

// StdDev returns the population standard deviation of the given values, or zero for an empty
// slice.
//
//	StdDev = Sqrt(1/N * Sum(Pow(value - mean, 2)))
//
// Example:
//
//	n := helper.StdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})
//	fmt.Println(n) // 2
func StdDev[T Float](values []T) T {
	if len(values) == 0 {
		return 0
	}

	meanValue := Mean(values)

	var sumSquaredDiff T

	for _, v := range values {
		diff := v - meanValue
		sumSquaredDiff += diff * diff
	}

	return T(math.Sqrt(float64(sumSquaredDiff / T(len(values)))))
}
