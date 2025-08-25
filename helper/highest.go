// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "slices"

// Highest returns a channel that emits the highest value
// within a sliding window of size w from the input channel c.
func Highest[T Number](c <-chan T, w int) <-chan T {
	return Window(c, func(s []T, i int) T {
		return slices.Max(s)
	}, w)
}
