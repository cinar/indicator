// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// MaxSinceWithContext returns a channel of T indicating since when
// (number of previous values) the maximum within the window of size w
// was seen. When the maximum occurs more than once, the most recent one
// counts. NaN values are ignored, and a window of only NaN values yields
// NaN.
func MaxSinceWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return extremeSinceWithContext(ctx, c, w, func(a, b T) bool { return a > b })
}

// MaxSince wraps MaxSinceWithContext for backwards compatibility.
//
// Deprecated: Use MaxSinceWithContext instead.
func MaxSince[T Number](c <-chan T, w int) <-chan T {
	return MaxSinceWithContext(context.Background(), c, w)
}
