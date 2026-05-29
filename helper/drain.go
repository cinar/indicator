// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Drain wraps DrainWithContext for backwards compatibility.
//
// Deprecated: Use DrainWithContext instead.
func Drain[T any](c <-chan T) {
	DrainWithContext(context.Background(), c)
}

// DrainWithContext drains the given channel with context support. It blocks the caller.
func DrainWithContext[T any](ctx context.Context, c <-chan T) {
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
