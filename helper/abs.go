// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"math"
)

// AbsWithContext calculates the absolute value of each value in a channel of type T.
//
// Example:
//
//	abs := helper.Abs(helper.SliceToChan([]int{-10, 20, -4, -5}))
//	fmt.Println(helper.ChanToSlice(abs)) // [10, 20, 4, 5]
func AbsWithContext[T Number](ctx context.Context, c <-chan T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		return T(math.Abs(float64(n)))
	})
}

// Abs wraps AbsWithContext for backwards compatibility.
//
// Deprecated: Use AbsWithContext instead.
func Abs[T Number](c <-chan T) <-chan T { return AbsWithContext(context.Background(), c) }
