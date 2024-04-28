// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Last takes a channel of values and returns a new channel containing the last N values.
func Last[T any](c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		ring := NewRing[T](count)

		for n := range c {
			ring.Put(n)
		}

		for !ring.IsEmpty() {
			n, _ := ring.Get()
			result <- n
		}
	}()

	return result
}
