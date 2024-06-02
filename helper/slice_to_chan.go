// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// SliceToChan converts a slice of float64 to a channel of float64.
//
// Example:
//
//	slice := []float64{2, 4, 6, 8}
//	c := helper.SliceToChan(slice)
//	fmt.Println(<- c)  // 2
//	fmt.Println(<- c)  // 4
//	fmt.Println(<- c)  // 6
//	fmt.Println(<- c)  // 8
func SliceToChan[T any](slice []T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for _, n := range slice {
			c <- n
		}
	}()

	return c
}
