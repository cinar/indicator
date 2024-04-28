// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Count generates a sequence of numbers starting with a specified value, from, and incrementing by one until
// the given other channel continues to produce values.
//
// Example:
//
//	other := make(chan int, 4)
//	other <- 1
//	other <- 1
//	other <- 1
//	other <- 1
//	close(other)
//
//	c := Count(0, other)
//
//	fmt.Println(<- s) // 1
//	fmt.Println(<- s) // 2
//	fmt.Println(<- s) // 3
//	fmt.Println(<- s) // 4
func Count[T Number, O any](from T, other <-chan O) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for i := from; ; i++ {
			_, ok := <-other
			if !ok {
				break
			}

			c <- i
		}
	}()

	return c
}
