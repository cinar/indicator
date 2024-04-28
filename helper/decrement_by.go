// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// DecrementBy decrements each element in the input channel by the
// specified decrement value and returns a new channel containing
// the decremented values.
//
// Example:
//
//	input := helper.SliceToChan([]int{1, 2, 3, 4})
//	substractOne := helper.DecrementBy(input, 1)
//	fmt.Println(helper.ChanToSlice(substractOne)) // [0, 1, 2, 3]
func DecrementBy[T Number](c <-chan T, d T) <-chan T {
	return Apply(c, func(n T) T {
		return n - d
	})
}
