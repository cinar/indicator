// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// ChanToSlice converts a channel of float64 to a slice of float64.
//
// Example:
//
//	c := make(chan int, 4)
//	c <- 1
//	c <- 2
//	c < -3
//	c <- 4
//	close(c)
//
//	fmt.Println(helper.ChanToSlice(c)) // [1, 2, 3, 4]
func ChanToSlice[T any](c <-chan T) []T {
	var slice []T

	for n := range c {
		slice = append(slice, n)
	}

	return slice
}
