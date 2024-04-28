// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// KeepNegatives processes a stream of float64 values, retaining negative
// values unchanged and replacing positive values with zero.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, 4, -5})
//	negatives := helper.KeepPositives(c)
//	fmt.Println(helper.ChanToSlice(negatives)) // [-10, 0, 0, -5]
func KeepNegatives[T Number](c <-chan T) <-chan T {
	return Apply(c, func(n T) T {
		if n < 0 {
			return n
		}

		return 0
	})
}
