// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Duplicate wraps DuplicateWithContext for backwards compatibility.
//
// Deprecated: Use DuplicateWithContext instead.
func Duplicate[T any](input <-chan T, count int) []<-chan T {
	return DuplicateWithContext(context.Background(), input, count)
}

// DuplicateWithContext duplicates a given receive-only channel by reading each value coming out of
// that channel and sending them on requested number of new output channels, supporting context cancellation.
func DuplicateWithContext[T any](ctx context.Context, input <-chan T, count int) []<-chan T {
	outputs := make([]chan T, count)
	result := make([]<-chan T, count)

	for i := range outputs {
		outputs[i] = make(chan T, cap(input))
		result[i] = outputs[i]
	}

	go func() {
		for _, output := range outputs {
			defer close(output)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-input:
				if !ok {
					return
				}
				for _, output := range outputs {
					select {
					case <-ctx.Done():
						return
					case output <- n:
					}
				}
			}
		}
	}()

	return result
}
