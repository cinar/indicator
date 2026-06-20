// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Buffered wraps BufferedWithContext for backwards compatibility.
//
// Deprecated: Use BufferedWithContext instead.
func Buffered[T any](c <-chan T, size int) <-chan T {
	return BufferedWithContext(context.Background(), c, size)
}

// BufferedWithContext takes a channel of any type and returns a new channel of the same type with
// a buffer of the specified size with context support.
func BufferedWithContext[T any](ctx context.Context, c <-chan T, size int) <-chan T {
	result := make(chan T, size)

	go PipeWithContext(ctx, c, result)

	return result
}
