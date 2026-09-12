// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "sync"

// ChanToSlices converts multiple channels of type T to slices of type T, draining all of the
// given channels concurrently on their own goroutines. This is the safe way to fully consume
// the unbuffered output channels returned by DuplicateWithContext (or any other multi-output
// producer), since draining them one at a time would block the producer forever once its
// internal buffer, if any, is exhausted.
//
// The returned slices are in the same order as the given channels.
//
// Example:
//
//	c1 := make(chan int, 4)
//	c2 := make(chan int, 4)
//
//	c1 <- 1
//	c2 <- 2
//	close(c1)
//	close(c2)
//
//	fmt.Println(helper.ChanToSlices(c1, c2)) // [[1] [2]]
func ChanToSlices[T any](chans ...<-chan T) [][]T {
	slices := make([][]T, len(chans))

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(chans))

	for i, c := range chans {
		go func(i int, c <-chan T) {
			defer waitGroup.Done()

			slices[i] = ChanToSlice(c)
		}(i, c)
	}

	waitGroup.Wait()

	return slices
}
