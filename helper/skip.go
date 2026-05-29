// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Skip wraps SkipWithContext for backwards compatibility.
//
// Deprecated: Use SkipWithContext instead.
func Skip[T any](c <-chan T, count int) <-chan T {
	return SkipWithContext(context.Background(), c, count)
}

// SkipWithContext skips the specified number of elements from the
// given channel of type T, supporting context cancellation.
func SkipWithContext[T any](ctx context.Context, c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				close(result)
				return
			case _, ok := <-c:
				if !ok {
					close(result)
					return
				}
			}
		}

		PipeWithContext(ctx, c, result)
	}()

	return result
}
