// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Gcd calculates the Greatest Common Divisor of the given numbers. It
// returns 0 for empty input, and treats negative values as their
// absolute value, since GCD is conventionally non-negative.
func Gcd(values ...int) int {
	if len(values) == 0 {
		return 0
	}

	gcd := abs(values[0])

	for i := 1; i < len(values); i++ {
		value := abs(values[i])

		for value > 0 {
			gcd, value = value, gcd%value
		}

		if gcd == 1 {
			break
		}
	}

	return gcd
}

// abs returns the absolute value of the given int.
func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
