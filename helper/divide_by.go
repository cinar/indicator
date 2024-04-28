// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// DivideBy divides each element in the input channel
// of float64 values by the given divider and returns a
// new channel containing the divided values.
//
// Example:
//
//	half := helper.DivideBy(helper.SliceToChan([]int{2, 4, 6, 8}), 2)
//	fmt.Println(helper.ChanToSlice(half)) // [1, 2, 3, 4]
func DivideBy[T Number](c <-chan T, d T) <-chan T {
	return Apply(c, func(n T) T {
		return n / d
	})
}
