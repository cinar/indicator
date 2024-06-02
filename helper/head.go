// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Head retrieves the specified number of elements
// from the given channel of float64 values and
// delivers them through a new channel.
//
// Example:
//
//	c := helper.SliceToChan([]int{2, 4, 6, 8})
//	actual := helper.Head(c, 2)
//	fmt.Println(helper.ChanToSlice(actual)) // [2, 4]
func Head[T Number](c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		for i := 0; i < count; i++ {
			n, ok := <-c
			if !ok {
				break
			}

			result <- n
		}
	}()

	return result
}
