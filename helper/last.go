// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Last wraps LastWithContext for backwards compatibility.
//
// Deprecated: Use LastWithContext instead.
func Last[T any](c <-chan T, count int) <-chan T {
	return LastWithContext(context.Background(), c, count)
}

// LastWithContext takes a channel of values and returns a new channel containing the last N values, supporting context cancellation.
func LastWithContext[T any](ctx context.Context, c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		ring := NewRing[T](count)

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					goto send
				}
				ring.Put(n)
			}
		}

	send:
		for !ring.IsEmpty() {
			n, _ := ring.Get()
			select {
			case <-ctx.Done():
				return
			case result <- n:
			}
		}
	}()

	return result
}
