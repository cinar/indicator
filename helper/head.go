// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Head wraps HeadWithContext for backwards compatibility.
//
// Deprecated: Use HeadWithContext instead.
func Head[T Number](c <-chan T, count int) <-chan T {
	return HeadWithContext(context.Background(), c, count)
}

// HeadWithContext retrieves the specified number of elements
// from the given channel of type T values and
// delivers them through a new channel, supporting context cancellation.
func HeadWithContext[T Number](ctx context.Context, c <-chan T, count int) <-chan T {
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
	}()

	return result
}
