// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Drain drains the given channel. It blocks the caller.
func Drain[T any](c <-chan T) {
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}
}
