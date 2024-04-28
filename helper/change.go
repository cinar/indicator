// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Change calculates the difference between the current value and the value N before.
//
// Example:
//
//	input := []int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4}
//	output := helper.Change(helper.SliceToChan(input), 2)
//	fmt.Println(helper.ChanToSlice(output)) // [4, 3, 3, -3, -7, -1, 2, 3]
func Change[T Number](c <-chan T, before int) <-chan T {
	cs := Duplicate(c, 2)
	cs[0] = Buffered(cs[0], before)
	cs[1] = Skip(cs[1], before)

	return Subtract(cs[1], cs[0])
}
