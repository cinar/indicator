// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Divide takes two channels of float64 and divides the values
// from the first channel with the values from the second one.
// It returns a new channel containing the results of
// the division.
//
// Example:
//
//	ac := helper.SliceToChan([]int{2, 4, 6, 8, 10})
//	bc := helper.SliceToChan([]int{2, 1, 3, 2, 5})
//
//	divison := helper.Divide(ac, bc)
//
//	fmt.Println(helper.ChanToSlice(division)) // [1, 4, 2, 4, 2]
func Divide[T Number](ac, bc <-chan T) <-chan T {
	return Operate(ac, bc, func(a, b T) T {
		return a / b
	})
}
