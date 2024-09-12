// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Gcd calculates the Greatest Common Divisor of the given numbers.
func Gcd(values ...int) int {
	gcd := values[0]

	for i := 1; i < len(values); i++ {
		value := values[i]

		for value > 0 {
			gcd, value = value, gcd%value
		}

		if gcd == 1 {
			break
		}
	}

	return gcd
}
