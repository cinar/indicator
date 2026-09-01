// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Lcm calculates the Least Common Multiple of the given numbers. It
// returns 0 for empty input, and returns 0 whenever 0 is among the
// values, matching the convention that LCM involving 0 is 0.
func Lcm(values ...int) int {
	if len(values) == 0 {
		return 0
	}

	lcm := values[0]

	for i := 1; i < len(values); i++ {
		if lcm == 0 || values[i] == 0 {
			return 0
		}

		// Divide before multiplying to reduce the risk of int overflow.
		lcm = values[i] / Gcd(values[i], lcm) * lcm
	}

	return lcm
}
