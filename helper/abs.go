// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"math"
)

// AbsWithContext calculates the absolute value of each value in a channel of type T.
//
// For an integer T, the absolute value is computed natively (negating
// the value directly) rather than through a float64 round-trip, so
// integers beyond float64's 53-bit exact-integer range (e.g., int64
// values beyond 2^53) remain exact.
//
// This does not change the behavior at the minimum value of a signed
// integer type (e.g., math.MinInt8, math.MinInt64), whose true
// absolute value does not fit in that same type. Negating the minimum
// value overflows and wraps back to the minimum value itself, which
// matches the result previously produced by the float64 round-trip.
// This is an inherent limitation of two's complement signed integers,
// not something either implementation can silently correct.
//
// Example:
//
//	abs := helper.Abs(helper.SliceToChan([]int{-10, 20, -4, -5}))
//	fmt.Println(helper.ChanToSlice(abs)) // [10, 20, 4, 5]
func AbsWithContext[T Number](ctx context.Context, c <-chan T) <-chan T {
	return ApplyWithContext(ctx, c, func(n T) T {
		switch any(n).(type) {
		case int, int8, int16, int32, int64:
			if n < 0 {
				return -n
			}

			return n
		}

		return T(math.Abs(float64(n)))
	})
}

// Abs wraps AbsWithContext for backwards compatibility.
//
// Deprecated: Use AbsWithContext instead.
func Abs[T Number](c <-chan T) <-chan T { return AbsWithContext(context.Background(), c) }
