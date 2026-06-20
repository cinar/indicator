// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"slices"
)

// LowestWithContext returns a channel that emits the lowest value
// within a sliding window of size w from the input channel c.
func LowestWithContext[T Number](ctx context.Context, c <-chan T, w int) <-chan T {
	return WindowWithContext(ctx, c, func(s []T, i int) T {
		return slices.Min(s)
	}, w)
}

// Lowest wraps LowestWithContext for backwards compatibility.
//
// Deprecated: Use LowestWithContext instead.
func Lowest[T Number](c <-chan T, w int) <-chan T {
	return LowestWithContext(context.Background(), c, w)
}
