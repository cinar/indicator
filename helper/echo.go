// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Echo wraps EchoWithContext for backwards compatibility.
//
// Deprecated: Use EchoWithContext instead.
func Echo[T any](input <-chan T, last, count int) <-chan T {
	return EchoWithContext(context.Background(), input, last, count)
}

// EchoWithContext takes a channel of numbers, repeats the specified count of numbers at the end by the specified count, supporting context cancellation.
func EchoWithContext[T any](ctx context.Context, input <-chan T, last, count int) <-chan T {
	output := make(chan T)
	memory := NewRing[T](last)

	go func() {
		defer close(output)

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-input:
				if !ok {
					goto repeat
				}
				memory.Put(n)
				select {
				case <-ctx.Done():
					return
				case output <- n:
				}
			}
		}

	repeat:
		for i := 0; i < count; i++ {
			for j := 0; j < last; j++ {
				select {
				case <-ctx.Done():
					return
				case output <- memory.At(j):
				}
			}
		}
	}()

	return output
}
