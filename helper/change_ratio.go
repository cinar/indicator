// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// ChangeRatioWithContext calculates the ratio change between the current
// value and the value N positions before.
//
// Example:
//
//	c := helper.SliceToChan([]float64{1, 2, 5, 5, 8, 2, 1, 1, 3, 4})
//	actual := helper.ChangeRatio(c, 2)
//	fmt.Println(helper.ChanToSlice(actual)) // [4, 1.5, 0.6, -0.6, -0.875, -0.5, 2, 3]
func ChangeRatioWithContext[T Number](ctx context.Context, c <-chan T, before int) <-chan T {
	cs := DuplicateWithContext(ctx, c, 2)
	cs[1] = BufferedWithContext(ctx, cs[1], before)
	return DivideWithContext(ctx, ChangeWithContext(ctx, cs[0], before), cs[1])
}

// ChangeRatio wraps ChangeRatioWithContext for backwards compatibility.
//
// Deprecated: Use ChangeRatioWithContext instead.
func ChangeRatio[T Number](c <-chan T, before int) <-chan T {
	return ChangeRatioWithContext(context.Background(), c, before)
}
