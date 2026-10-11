// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// SafeDivide divides the given numerator by the given denominator, and
// returns the given fallback instead of dividing by zero when the
// denominator is zero.
//
// Several indicators guard a division whose denominator can be
// legitimately zero (e.g., a flat high-low range, or a flat window with
// no true range at all), and return a neutral value instead of letting
// the division produce NaN or Inf. This centralizes that guard so each
// indicator only needs to supply its own neutral fallback.
//
// Example:
//
//	n := helper.SafeDivide(4.0, 2.0, 0.5)
//	fmt.Println(n) // 2
//
//	n = helper.SafeDivide(4.0, 0.0, 0.5)
//	fmt.Println(n) // 0.5
func SafeDivide[T Number](numerator, denominator, fallback T) T {
	if denominator == 0 {
		return fallback
	}

	return numerator / denominator
}
