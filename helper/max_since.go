// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"slices"
)

// MaxSinceWithContext returns a channel of T indicating since when
// (number of previous values) the respective value was the maximum
// within the window of size w.
func MaxSinceWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return WindowWithContext(ctx, c, func(w []T, i int) T {
		since := 0
		found := false
		m := slices.Max(w)
		SlicesReverse(w, i, func(n T) bool {
			if found && n < m {
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

// MaxSince wraps MaxSinceWithContext for backwards compatibility.
//
// Deprecated: Use MaxSinceWithContext instead.
func MaxSince[T Number](c <-chan T, w int) <-chan T {
	return MaxSinceWithContext(context.Background(), c, w)
}
