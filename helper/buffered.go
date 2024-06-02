// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Buffered takes a channel of any type and returns a new channel of the same type with
// a buffer of the specified size. This allows the original channel to continue sending
// data even if the receiving end is temporarily unavailable.
//
// Example:
func Buffered[T any](c <-chan T, size int) <-chan T {
	result := make(chan T, size)

	go Pipe(c, result)

	return result
}
