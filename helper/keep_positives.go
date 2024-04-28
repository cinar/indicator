// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// KeepPositives processes a stream of float64 values, retaining positive
// values unchanged and replacing negative values with zero.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, 4, -5})
//	positives := helper.KeepPositives(c)
//	fmt.Println(helper.ChanToSlice(positives)) // [0, 20, 4, 0]
func KeepPositives[T Number](c <-chan T) <-chan T {
	return Apply(c, func(n T) T {
		if n > 0 {
			return n
		}

		return 0
	})
}
