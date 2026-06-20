// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// MultiplyWithContext takes two channels of type T and multiples the values
// from the first channel with the values from the second channel.
// It returns a new channel containing the results of
// the multiplication.
//
// Example:
//
//	ac := helper.SliceToChan([]int{1, 4, 2, 4, 2})
//	bc := helper.SliceToChan([]int{2, 1, 3, 2, 5})
//
//	multiplication := helper.Multiply(ac, bc)
//
//	fmt.Println(helper.ChanToSlice(multiplication)) // [2, 4, 6, 8, 10]
func MultiplyWithContext[T Number](ctx context.Context, ac, bc <-chan T) <-chan T {
	return OperateWithContext(ctx, ac, bc, func(a, b T) T {
		return a * b
	})
}

// Multiply wraps MultiplyWithContext for backwards compatibility.
//
// Deprecated: Use MultiplyWithContext instead.
func Multiply[T Number](ac, bc <-chan T) <-chan T {
	return MultiplyWithContext(context.Background(), ac, bc)
}
