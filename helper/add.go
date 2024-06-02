// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Add adds each pair of values from the two input channels of float64
// and returns a new channel containing the sums.
//
// Example:
//
//	ac := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//
//	actual := helper.ChanToSlice(helper.Add(ac, bc))
//
//	fmt.Println(actual) // [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]
func Add[T Number](ac, bc <-chan T) <-chan T {
	return Operate(ac, bc, func(a, b T) T {
		return a + b
	})
}
