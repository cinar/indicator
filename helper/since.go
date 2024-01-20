// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper

// Since counts the number of periods since the last change of
// value in a channel of numbers.
func Since[T Number](c <-chan T) <-chan T {
	first := true
	last := T(0)
	count := T(0)

	return Apply(c, func(n T) T {
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
