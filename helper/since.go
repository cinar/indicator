// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// SinceWithContext counts the number of periods since the last change of
// value in a channel of numbers, supporting context cancellation.
func SinceWithContext[T comparable, R Number](ctx context.Context, c <-chan T) <-chan R {
	first := true

	var last T
	var count R

	return MapWithContext(ctx, c, func(n T) R {
		if first || last != n {
			first = false
			last = n
			count = 0
		} else {
			count++
		}

		return count
	})
}

// Since wraps SinceWithContext for backwards compatibility.
//
// Deprecated: Use SinceWithContext instead.
func Since[T comparable, R Number](c <-chan T) <-chan R {
	return SinceWithContext[T, R](context.Background(), c)
}
