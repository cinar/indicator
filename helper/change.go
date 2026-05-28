// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// ChangeWithContext calculates the difference between the current value and the value N before.
//
// Example:
//
//	input := []int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4}
//	output := helper.Change(helper.SliceToChan(input), 2)
//	fmt.Println(helper.ChanToSlice(output)) // [4, 3, 3, -3, -7, -1, 2, 3]
func ChangeWithContext[T Number](ctx context.Context, c <-chan T, before int) <-chan T {
	cs := DuplicateWithContext(ctx, c, 2)
	cs[0] = BufferedWithContext(ctx, cs[0], before)
	cs[1] = SkipWithContext(ctx, cs[1], before)

	return Subtract(cs[1], cs[0])
}

// Change wraps ChangeWithContext for backwards compatibility.
//
// Deprecated: Use ChangeWithContext instead.
func Change[T Number](c <-chan T, before int) <-chan T {
	return ChangeWithContext(context.Background(), c, before)
}
