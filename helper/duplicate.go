// Copyright (c) 2021-2026 The Indicator Authors.
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
//
// The returned output channels are unbuffered (regardless of the input channel's buffering), so all of
// them must be actively read concurrently. Fully draining one output channel before starting to read
// another will block the producer goroutine forever, since it sends to every output channel for each
// value before moving on to the next. Use ChanToSlices to safely drain all of the returned channels
// concurrently, or otherwise ensure each channel is read from its own goroutine.
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
