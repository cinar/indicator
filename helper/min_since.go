// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"slices"
)

// MinSince returns a channel of T indicating since when
// (number of previous values) the respective value was the minimum.
func MinSince[T Number](c <-chan T, w int) <-chan T {
	return Window(c, func(w []T, i int) T {
		since := 0
		found := false
		m := slices.Min(w)
		SlicesReverse(w, i, func(n T) bool {
			if found && n > m {
				return false
			}
			since++
			if n == m {
				found = true
			}
			return true
		})
		return T(since - 1)
	}, w)
}
