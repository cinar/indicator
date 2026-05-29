// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"sync"
)

// Waitable wraps WaitableWithContext for backwards compatibility.
//
// Deprecated: Use WaitableWithContext instead.
func Waitable[T any](wg *sync.WaitGroup, c <-chan T) <-chan T {
	return WaitableWithContext(context.Background(), wg, c)
}

// WaitableWithContext increments the wait group before reading from the channel
// and signals completion when the channel is closed, supporting context cancellation.
func WaitableWithContext[T any](ctx context.Context, wg *sync.WaitGroup, c <-chan T) <-chan T {
	result := make(chan T, cap(c))

	wg.Add(1)

	go func() {
		defer close(result)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case result <- n:
				}
			}
		}
	}()

	return result
}
