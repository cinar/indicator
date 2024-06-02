// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Since counts the number of periods since the last change of
// value in a channel of numbers.
func Since[T comparable, R Number](c <-chan T) <-chan R {
	first := true

	var last T
	var count R

	return Map(c, func(n T) R {
		if first || last != n {
			first = false
			last = n
			count = 0
		} else {
			count++
		}

		return count
	})
}
