// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// IsRising returns 1 if the current value is strictly greater than the value
// n periods ago, and 0 otherwise.
//
// Example:
//
//	input := []int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4}
//	output := helper.IsRising(helper.SliceToChan(input), 2)
//	fmt.Println(helper.ChanToSlice(output)) // [1, 1, 1, 0, 0, 0, 1, 1]
func IsRising[T Number](c <-chan T, period int) <-chan T {
	cs := Duplicate(c, 2)
	cs[0] = Buffered(cs[0], period)
	cs[1] = Skip(cs[1], period)

	return Operate(cs[1], cs[0], func(current, previous T) T {
		if current > previous {
			return 1
		}

		return 0
	})
}
