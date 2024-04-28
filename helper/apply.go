// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Apply applies the given transformation function to each element in the
// input channel and returns a new channel containing the transformed
// values. The transformation function takes a float64 value as input
// and returns a float64 value as output.
//
// Example:
//
//	timesTwo := helper.Apply(c, func(n int) int {
//		return n * 2
//	})
func Apply[T Number](c <-chan T, f func(T) T) <-chan T {
	ac := make(chan T)

	go func() {
		defer close(ac)

		for n := range c {
			ac <- f(n)
		}
	}()

	return ac
}
