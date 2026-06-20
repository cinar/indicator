// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// KeepPositivesWithContext processes a stream of type T values, retaining positive
// values unchanged and replacing negative values with zero.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, 4, -5})
//	positives := helper.KeepPositives(c)
//	fmt.Println(helper.ChanToSlice(positives)) // [0, 20, 4, 0]
func KeepPositivesWithContext[T Number](ctx context.Context, c <-chan T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		if n > 0 {
			return n
		}

		return 0
	})
}

// KeepPositives wraps KeepPositivesWithContext for backwards compatibility.
//
// Deprecated: Use KeepPositivesWithContext instead.
func KeepPositives[T Number](c <-chan T) <-chan T {
	return KeepPositivesWithContext(context.Background(), c)
}
