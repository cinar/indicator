// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// RoundDigits takes a channel of float64 numbers and rounds them to d
// decimal places.
//
// Example:
//
//	c := helper.SliceToChan([]float64{10.1234, 5.678, 6.78, 8.91011})
//	rounded := helper.RoundDigits(c, 2)
//	fmt.Println(helper.ChanToSlice(rounded)) // [10.12, 5.68, 6.78, 8.91]
func RoundDigits[T Number](c <-chan T, d int) <-chan T {
	return Apply(c, func(n T) T {
		return RoundDigit(n, d)
	})
}
