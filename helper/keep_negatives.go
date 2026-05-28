// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// KeepNegativesWithContext processes a stream of type T values, retaining negative
// values unchanged and replacing positive values with zero.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, 4, -5})
//	negatives := helper.KeepPositives(c)
//	fmt.Println(helper.ChanToSlice(negatives)) // [-10, 0, 0, -5]
func KeepNegativesWithContext[T Number](ctx context.Context, c <-chan T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		if n < 0 {
			return n
		}

		return 0
	})
}

// KeepNegatives wraps KeepNegativesWithContext for backwards compatibility.
//
// Deprecated: Use KeepNegativesWithContext instead.
func KeepNegatives[T Number](c <-chan T) <-chan T {
	return KeepNegativesWithContext(context.Background(), c)
}
