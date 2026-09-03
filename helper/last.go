// Copyright (c) 2021-2026 The Indicator Authors.
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
//
// A count of zero or less means there is nothing to keep: the input
// channel is still drained, so upstream stages are not blocked, but no
// values are emitted. (Unlike EchoWithContext, clamping a ring to a
// minimum size of 1 would not be a safe stand-in here, since it would
// wrongly surface one value instead of none, so this is guarded
// explicitly rather than delegated to NewRing.)
func LastWithContext[T any](ctx context.Context, c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		if count <= 0 {
			for {
				select {
				case <-ctx.Done():
					return
				case _, ok := <-c:
					if !ok {
						return
					}
				}
			}
		}

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
