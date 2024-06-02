// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Map applies the given transformation function to each element in the
// input channel and returns a new channel containing the transformed
// values. The transformation function takes a float64 value as input
// and returns a float64 value as output.
//
// Example:
//
//	timesTwo := helper.Map(c, func(n int) int {
//		return n * 2
//	})
func Map[F, T any](c <-chan F, f func(F) T) <-chan T {
	mc := make(chan T)

	go func() {
		defer close(mc)

		for n := range c {
			mc <- f(n)
		}
	}()

	return mc
}
