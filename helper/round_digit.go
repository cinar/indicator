// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "math"

// RoundDigit rounds the given float64 number to d decimal places.
//
// For an integer T, rounding to decimal places is a no-op, since an
// integer has no fractional component to round away. In that case, n
// is returned unchanged, which also avoids routing the value through
// a lossy float64 round-trip for integers beyond float64's 53-bit
// exact-integer range (e.g., int64 values beyond 2^53).
//
// Example:
//
//	n := helper.RoundDigit(10.1234, 2)
//	fmt.Println(n) // 10.12
func RoundDigit[T Number](n T, d int) T {
	switch any(n).(type) {
	case int, int8, int16, int32, int64:
		return n
	}

	m := math.Pow(10, float64(d))
	return T(math.Round(float64(n)*m) / m)
}
