// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// IncrementBy increments each element in the input channel by the
// specified increment value and returns a new channel containing
// the incremented values.
//
// Example:
//
//	input := []int{1, 2, 3, 4}
//	actual := helper.IncrementBy(helper.SliceToChan(input), 1)
//	fmt.Println(helper.ChanToSlice(actual)) // [2, 3, 4, 5]
func IncrementBy[T Number](c <-chan T, i T) <-chan T {
	return Apply(c, func(n T) T {
		return n + i
	})
}
