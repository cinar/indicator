// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Map wraps MapWithContext for backwards compatibility.
//
// Deprecated: Use MapWithContext instead.
func Map[F, T any](c <-chan F, f func(F) T) <-chan T {
	return MapWithContext(context.Background(), c, f)
}

// MapWithContext applies the given transformation function to each element in the
// input channel and returns a new channel containing the transformed
// values, supporting context cancellation.
func MapWithContext[F, T any](ctx context.Context, c <-chan F, f func(F) T) <-chan T {
	mc := make(chan T)

	go func() {
		defer close(mc)

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case mc <- f(n):
				}
			}
		}
	}()

	return mc
}
