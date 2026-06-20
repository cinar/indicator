// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Pipe wraps PipeWithContext for backwards compatibility.
//
// Deprecated: Use PipeWithContext instead.
func Pipe[T any](f <-chan T, t chan<- T) {
	PipeWithContext(context.Background(), f, t)
}

// PipeWithContext copies all elements from the input channel into the output channel
// with context support.
func PipeWithContext[T any](ctx context.Context, f <-chan T, t chan<- T) {
	defer close(t)
	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-f:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case t <- n:
			}
		}
	}
}
