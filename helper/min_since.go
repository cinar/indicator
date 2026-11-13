// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// MinSinceWithContext returns a channel of T indicating since when
// (number of previous values) the minimum within the window of size w
// was seen. When the minimum occurs more than once, the most recent one
// counts. NaN values are ignored, and a window of only NaN values yields
// NaN.
func MinSinceWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return extremeSinceWithContext(ctx, c, w, func(a, b T) bool { return a < b })
}

// MinSince wraps MinSinceWithContext for backwards compatibility.
//
// Deprecated: Use MinSinceWithContext instead.
func MinSince[T Number](c <-chan T, w int) <-chan T {
	return MinSinceWithContext(context.Background(), c, w)
}
