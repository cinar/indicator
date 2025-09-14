// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "slices"

// Lowest returns a channel that emits the lowest value
// within a sliding window of size w from the input channel c.
func Lowest[T Number](c <-chan T, w int) <-chan T {
	return Window(c, func(s []T, i int) T {
		return slices.Min(s)
	}, w)
}
