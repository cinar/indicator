// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// IsRisingWithContext takes a channel of type T values and returns 1 if the current
// value is strictly greater than the value the given period before, and 0 otherwise,
// supporting context cancellation.
//
// Example:
//
//	input := []int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4}
//	output := helper.IsRising(helper.SliceToChan(input), 1)
//	fmt.Println(helper.ChanToSlice(output)) // [1, 1, 0, 1, 0, 0, 0, 1, 1]
func IsRisingWithContext[T Number](ctx context.Context, c <-chan T, period int) <-chan T {
	cs := DuplicateWithContext(ctx, c, 2)
	cs[0] = BufferedWithContext(ctx, cs[0], period)
	cs[1] = SkipWithContext(ctx, cs[1], period)

	return OperateWithContext(ctx, cs[1], cs[0], func(current, before T) T {
		if current > before {
			return 1
		}

		return 0
	})
}

// IsRising wraps IsRisingWithContext for backwards compatibility.
//
// Deprecated: Use IsRisingWithContext instead.
func IsRising[T Number](c <-chan T, period int) <-chan T {
	return IsRisingWithContext(context.Background(), c, period)
}
