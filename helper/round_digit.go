// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "math"

// RoundDigit rounds the given float64 number to d decimal places.
//
// Example:
//
//	n := helper.RoundDigit(10.1234, 2)
//	fmt.Println(n) // 10.12
func RoundDigit[T Number](n T, d int) T {
	m := math.Pow(10, float64(d))
	return T(math.Round(float64(n)*m) / m)
}
