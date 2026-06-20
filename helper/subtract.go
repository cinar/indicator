// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
)

// SubtractWithContext takes two channels of type T and subtracts the values
// from the second channel from the first one. It returns a new
// channel containing the results of the subtractions.
//
// Example:
//
//	ac := helper.SliceToChan([]int{2, 4, 6, 8, 10})
//	bc := helper.SliceToChan([]int{1, 2, 3, 4, 5})
//	actual := helper.Subtract(ac, bc)
//	fmt.Println(helper.ChanToSlice(actual)) // [1, 2, 3, 4, 5]
func SubtractWithContext[T Number](ctx context.Context, ac, bc <-chan T) <-chan T {
	return OperateWithContext(ctx, ac, bc, func(a, b T) T {
		return a - b
	})
}

// Subtract wraps SubtractWithContext for backwards compatibility.
//
// Deprecated: Use SubtractWithContext instead.
func Subtract[T Number](ac, bc <-chan T) <-chan T {
	return SubtractWithContext(context.Background(), ac, bc)
}
