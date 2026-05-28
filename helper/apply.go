// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Apply wraps ApplyWithContext for backwards compatibility.
//
// Deprecated: Use ApplyWithContext instead.
func Apply[T Number](c <-chan T, f func(T) T) <-chan T {
	return ApplyWithContext(context.Background(), c, f)
}

// ApplyWithContext applies the given transformation function to each element in the
// input channel and returns a new channel containing the transformed
// values with context support.
func ApplyWithContext[T Number](ctx context.Context, c <-chan T, f func(T) T) <-chan T {
	ac := make(chan T)

	go func() {
		defer close(ac)

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
				case ac <- f(n):
				}
			}
		}
	}()

	return ac
}
