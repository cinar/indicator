// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "sync"

// Waitable increments the wait group before reading from the channel
// and signals completion when the channel is closed.
func Waitable[T any](wg *sync.WaitGroup, c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	wg.Add(1)

	go func() {
		defer close(result)
		defer wg.Done()

		for n := range c {
			result <- n
		}
	}()

	return result
}
