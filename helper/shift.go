// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Shift takes a channel of numbers, shifts them to the right by the specified count,
// and fills in any missing values with the provided fill value.
//
// Example:
//
//	input := helper.SliceToChan([]int{2, 4, 6, 8})
//	output := helper.ChanToSlice(input, 4, 0))
//	fmt.Println(helper.ChanToSlice(output)) // [0, 0, 0, 0, 2, 4, 6, 8]
func Shift[T any](c <-chan T, count int, fill T) <-chan T {
	result := make(chan T, cap(c)+count)

	go func() {
		for i := 0; i < count; i++ {
			result <- fill
		}

		Pipe(c, result)
	}()

	return result
}
