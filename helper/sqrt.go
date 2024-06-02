// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "math"

// Sqrt calculates the square root of each value in a channel of float64.
//
// Example:
//
//	c := helper.SliceToChan([]int{9, 81, 16, 100})
//	sqrt := helper.Sqrt(c)
//	fmt.Println(helper.ChanToSlice(sqrt)) // [3, 9, 4, 10]
func Sqrt[T Number](c <-chan T) <-chan T {
	return Apply(c, func(n T) T {
		return T(math.Sqrt(float64(n)))
	})
}
