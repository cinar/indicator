// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Sign takes a channel of float64 values and returns their signs
// as -1 for negative, 0 for zero, and 1 for positive.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, -4, 0})
//	sign := helper.Sign(c)
//	fmt.Println(helper.ChanToSlice(sign)) // [-1, 1, -1, 0]
func Sign[T Number](c <-chan T) <-chan T {
	return Apply(c, func(n T) T {
		if n > 0 {
			return 1
		} else if n < 0 {
			return -1
		}

		return 0
	})
}
