// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Shift wraps ShiftWithContext for backwards compatibility.
//
// Deprecated: Use ShiftWithContext instead.
func Shift[T any](c <-chan T, count int, fill T) <-chan T {
	return ShiftWithContext(context.Background(), c, count, fill)
}

// ShiftWithContext takes a channel of numbers, shifts them to the right by the specified count,
// and fills in any missing values with the provided fill value, supporting context cancellation.
func ShiftWithContext[T any](ctx context.Context, c <-chan T, count int, fill T) <-chan T {
	result := make(chan T, cap(c)+count)

	go func() {
		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				close(result)
				return
			case result <- fill:
			}
		}

		PipeWithContext(ctx, c, result)
	}()

	return result
}
