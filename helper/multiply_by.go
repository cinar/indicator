// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// MultiplyBy multiplies each element in the input channel
// of float64 values by the given multiplier and returns a
// new channel containing the multiplied values.
//
// Example:
//
//	c := helper.SliceToChan([]int{1, 2, 3, 4})
//	twoTimes := helper.MultiplyBy(c, 2)
//	fmt.Println(helper.ChanToSlice(twoTimes)) // [2, 4, 6, 8]
func MultiplyBy[T Number](c <-chan T, m T) <-chan T {
	return Apply(c, func(n T) T {
		return n * m
	})
}
