// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// First takes a channel of values and returns a new channel containing the first N values.
func First[T any](c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		for i := 0; i < count; i++ {
			n, ok := <-c
			if !ok {
				break
			}

			result <- n
		}

		close(result)

		Drain(c)
	}()

	return result
}
