// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "math"

// Pow takes a channel of float64 values and returns the element-wise
// base-value exponential of y.
//
// Example:
//
//	c := helper.SliceToChan([]int{2, 3, 5, 10})
//	squared := helper.Pow(c, 2)
//	fmt.Println(helper.ChanToSlice(squared)) // [4, 9, 25, 100]
func Pow[T Number](c <-chan T, y T) <-chan T {
	return Apply(c, func(n T) T {
		return T(math.Pow(float64(n), float64(y)))
	})
}
