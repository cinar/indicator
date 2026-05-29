// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// DivideByWithContext divides each element in the input channel
// of type T values by the given divider and returns a
// new channel containing the divided values.
//
// Example:
//
//	half := helper.DivideBy(helper.SliceToChan([]int{2, 4, 6, 8}), 2)
//	fmt.Println(helper.ChanToSlice(half)) // [1, 2, 3, 4]
func DivideByWithContext[T Number](ctx context.Context, c <-chan T, d T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		return n / d
	})
}

// DivideBy wraps DivideByWithContext for backwards compatibility.
//
// Deprecated: Use DivideByWithContext instead.
func DivideBy[T Number](c <-chan T, d T) <-chan T {
	return DivideByWithContext(context.Background(), c, d)
}
