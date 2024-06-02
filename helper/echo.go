// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Echo takes a channel of numbers, repeats the specified count of numbers at the end by the specified count.
//
// Example:
//
//	input := helper.SliceToChan([]int{2, 4, 6, 8})
//	output := helper.Echo(input, 2, 4))
//	fmt.Println(helper.ChanToSlice(output)) // [2, 4, 6, 8, 6, 8, 6, 8, 6, 8, 6, 8]
func Echo[T any](input <-chan T, last, count int) <-chan T {
	output := make(chan T)
	memory := NewRing[T](last)

	go func() {
		defer close(output)

		for n := range input {
			memory.Put(n)
			output <- n
		}

		for i := 0; i < count; i++ {
			for j := 0; j < last; j++ {
				output <- memory.At(j)
			}
		}
	}()

	return output
}
