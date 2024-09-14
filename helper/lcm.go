// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Lcm calculates the Least Common Multiple of the given numbers.
func Lcm(values ...int) int {
	lcm := values[0]

	for i := 1; i < len(values); i++ {
		lcm = (values[i] * lcm) / Gcd(values[i], lcm)
	}

	return lcm
}
