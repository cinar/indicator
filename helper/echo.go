// Copyright (c) 2021-2026 The Indicator Authors.
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
//
// A last of zero or less means there is nothing to replay: the repeat
// phase emits no values regardless of count, and only the original
// input is forwarded (NewRing clamps its backing ring to a minimum
// size of 1 in this case, but that ring's contents are never read).
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
				v, ok := memory.At(j)
				if !ok {
					break
				}

				select {
				case <-ctx.Done():
					return
				case output <- v:
				}
			}
		}
	}()

	return output
}
