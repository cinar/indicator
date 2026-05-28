// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Filter wraps FilterWithContext for backwards compatibility.
//
// Deprecated: Use FilterWithContext instead.
func Filter[T any](c <-chan T, p func(T) bool) <-chan T {
	return FilterWithContext(context.Background(), c, p)
}

// FilterWithContext filters the items from the input channel based on the
// provided predicate function, supporting context cancellation.
func FilterWithContext[T any](ctx context.Context, c <-chan T, p func(T) bool) <-chan T {
	fc := make(chan T)

	go func() {
		defer close(fc)

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				if p(n) {
					select {
					case <-ctx.Done():
						return
					case fc <- n:
					}
				}
			}
		}
	}()

	return fc
}
