// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Seq generates a sequence of numbers starting with a specified value,
// from, and incrementing by a specified amount, increment, until a
// specified value, to, is reached or exceeded. The sequence includes
// both from and to.
//
// Example:
//
//	s := Seq(1, 5, 1)
//	defer close(s)
//
//	fmt.Println(<- s) // 1
//	fmt.Println(<- s) // 2
//	fmt.Println(<- s) // 3
//	fmt.Println(<- s) // 4
func Seq[T Number](from, to, increment T) <-chan T {
	c := make(chan T)

	go func() {
		for i := from; i < to; i += increment {
			c <- i
		}

		close(c)
	}()

	return c
}
