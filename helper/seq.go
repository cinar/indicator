// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Seq wraps SeqWithContext for backwards compatibility.
//
// Deprecated: Use SeqWithContext instead.
func Seq[T Number](from, to, increment T) <-chan T {
	return SeqWithContext(context.Background(), from, to, increment)
}

// SeqWithContext generates a sequence of numbers, supporting context cancellation.
func SeqWithContext[T Number](ctx context.Context, from, to, increment T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)
		for i := from; i < to; i += increment {
			select {
			case <-ctx.Done():
				return
			case c <- i:
			}
		}
	}()

	return c
}
