// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Count wraps CountWithContext for backwards compatibility.
//
// Deprecated: Use CountWithContext instead.
func Count[T Number, O any](from T, other <-chan O) <-chan T {
	return CountWithContext(context.Background(), from, other)
}

// CountWithContext generates a sequence of numbers starting with a specified value, from, and incrementing by one until
// the given other channel continues to produce values, supporting context cancellation.
func CountWithContext[T Number, O any](ctx context.Context, from T, other <-chan O) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for i := from; ; i++ {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-other:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case c <- i:
				}
			}
		}
	}()

	return c
}
