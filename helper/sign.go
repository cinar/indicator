// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// SignWithContext takes a channel of type T values and returns their signs
// as -1 for negative, 0 for zero, and 1 for positive.
//
// Example:
//
//	c := helper.SliceToChan([]int{-10, 20, -4, 0})
//	sign := helper.Sign(c)
//	fmt.Println(helper.ChanToSlice(sign)) // [-1, 1, -1, 0]
func SignWithContext[T Number](ctx context.Context, c <-chan T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		if n > 0 {
			return 1
		} else if n < 0 {
			return -1
		}

		return 0
	})
}

// Sign wraps SignWithContext for backwards compatibility.
//
// Deprecated: Use SignWithContext instead.
func Sign[T Number](c <-chan T) <-chan T { return SignWithContext(context.Background(), c) }
