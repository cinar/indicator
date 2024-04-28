// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Subtract takes two channels of float64 and subtracts the values
// from the second channel from the first one. It returns a new
// channel containing the results of the subtractions.
//
// Example:
//
//	ac := helper.SliceToChan([]int{2, 4, 6, 8, 10})
//	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
//	actual := helper.Subtract(ac, bc)
//	fmt.Println(helper.ChanToSlice(actual)) // [1, 2, 3, 4, 5]
func Subtract[T Number](ac, bc <-chan T) <-chan T {
	return Operate(ac, bc, func(a, b T) T {
		return a - b
	})
}
