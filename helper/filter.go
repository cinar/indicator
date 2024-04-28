// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Filter filters the items from the input channel based on the
// provided predicate function. The predicate function takes a
// float64 value as input and returns a boolean value indicating
// whether the value should be included in the output channel.
//
// Example:
//
//	even := helper.Filter(c, func(n int) bool {
//	  return n%2 == 0
//	})
func Filter[T any](c <-chan T, p func(T) bool) <-chan T {
	fc := make(chan T)

	go func() {
		for n := range c {
			if p(n) {
				fc <- n
			}
		}

		close(fc)
	}()

	return fc
}
