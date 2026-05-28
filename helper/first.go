// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// First wraps FirstWithContext for backwards compatibility.
//
// Deprecated: Use FirstWithContext instead.
func First[T any](c <-chan T, count int) <-chan T {
	return FirstWithContext(context.Background(), c, count)
}

// FirstWithContext takes a channel of values and returns a new channel containing the first N values, supporting context cancellation.
func FirstWithContext[T any](ctx context.Context, c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		for i := 0; i < count; i++ {
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
				case result <- n:
				}
			}
		}

		DrainWithContext(ctx, c)
	}()

	return result
}
