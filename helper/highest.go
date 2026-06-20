// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"slices"
)

// HighestWithContext returns a channel that emits the highest value
// within a sliding window of size w from the input channel c.
func HighestWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return WindowWithContext(ctx, c, func(s []T, i int) T {
		return slices.Max(s)
	}, w)
}

// Highest wraps HighestWithContext for backwards compatibility.
//
// Deprecated: Use HighestWithContext instead.
func Highest[T Number](c <-chan T, w int) <-chan T {
	return HighestWithContext(context.Background(), c, w)
}
