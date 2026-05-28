// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// MapWithPrevious wraps MapWithPreviousWithContext for backwards compatibility.
//
// Deprecated: Use MapWithPreviousWithContext instead.
func MapWithPrevious[F, T any](c <-chan F, f func(T, F) T, previous T) <-chan T {
	return MapWithPreviousWithContext(context.Background(), c, f, previous)
}

// MapWithPreviousWithContext applies a transformation function to each element in an input channel, creating a new channel
// with the transformed values, supporting context cancellation.
func MapWithPreviousWithContext[F, T any](ctx context.Context, c <-chan F, f func(T, F) T, previous T) <-chan T {
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
				previous = f(previous, n)
				select {
				case <-ctx.Done():
					return
				case mc <- previous:
				}
			}
		}
	}()

	return mc
}
