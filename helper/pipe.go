// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Pipe function takes an input channel and an output channel and copies
// all elements from the input channel into the output channel.
//
// Example:
//
//	input := helper.SliceToChan([]int{2, 4, 6, 8})
//	output := make(chan int)
//	helper.Pipe(input, output)
//	fmt.println(helper.ChanToSlice(output)) // [2, 4, 6, 8]
func Pipe[T any](f <-chan T, t chan<- T) {
	defer close(t)
	for n := range f {
		t <- n
	}
}
