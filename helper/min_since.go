// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"slices"
)

// MinSinceWithContext returns a channel of T indicating since when
// (number of previous values) the respective value was the minimum.
func MinSinceWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return WindowWithContext(ctx, c, func(w []T, i int) T {
		since := 0
		found := false
		m := slices.Min(w)
		SlicesReverse(w, i, func(n T) bool {
			if found && n > m {
				return false
			}
			since++
			if n == m {
				found = true
			}
			return true
		})
		return T(since - 1)
	}, w)
}

// MinSince wraps MinSinceWithContext for backwards compatibility.
//
// Deprecated: Use MinSinceWithContext instead.
func MinSince[T Number](c <-chan T, w int) <-chan T {
	return MinSinceWithContext(context.Background(), c, w)
}
