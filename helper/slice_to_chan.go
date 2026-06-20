// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// SliceToChan wraps SliceToChanWithContext for backwards compatibility.
//
// Deprecated: Use SliceToChanWithContext instead.
func SliceToChan[T any](slice []T) <-chan T {
	return SliceToChanWithContext(context.Background(), slice)
}

// SliceToChanWithContext converts a slice of type T to a channel of type T, supporting context cancellation.
func SliceToChanWithContext[T any](ctx context.Context, slice []T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for _, n := range slice {
			select {
			case <-ctx.Done():
				return
			case c <- n:
			}
		}
	}()

	return c
}
