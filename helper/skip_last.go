// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// SkipLast wraps SkipLastWithContext for backwards compatibility.
//
// Deprecated: Use SkipLastWithContext instead.
func SkipLast[T any](c <-chan T, count int) <-chan T {
	return SkipLastWithContext(context.Background(), c, count)
}

// SkipLastWithContext skips the specified number of elements
// from the end of the given channel, supporting context cancellation.
//
// A negative count is treated as zero: nothing is skipped, and every
// input value is forwarded unchanged, matching SkipWithContext's
// negative-count convention.
func SkipLastWithContext[T any](ctx context.Context, c <-chan T, count int) <-chan T {
	if count < 0 {
		count = 0
	}

	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		buf := make([]T, 0, count)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				buf = append(buf, v)
				if len(buf) > count {
					select {
					case <-ctx.Done():
						return
					case result <- buf[0]:
					}
					buf = buf[1:]
				}
			}
		}
	}()

	return result
}
