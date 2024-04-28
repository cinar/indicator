// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// MapWithPrevious applies a transformation function to each element in an input channel, creating a new channel
// with the transformed values. It maintains a "memory" of the previous result, allowing the transformation
// function to consider both the current element and the outcome of the previous transformation. This
// enables functions that rely on accumulated state or sequential dependencies between elements.
//
// Example:
//
//	sum := helper.MapWithPrevious(c, func(p, c int) int {
//		return p + c
//	}, 0)
func MapWithPrevious[F, T any](c <-chan F, f func(T, F) T, previous T) <-chan T {
	mc := make(chan T)

	go func() {
		defer close(mc)

		for n := range c {
			previous = f(previous, n)
			mc <- previous
		}
	}()

	return mc
}
