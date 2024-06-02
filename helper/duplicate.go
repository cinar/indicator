// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Duplicate duplicates a given receive-only channel by reading each value coming out of
// that channel and sending them on requested number of new output channels.
//
// Example:
//
//	expected := helper.SliceToChan([]float64{-10, 20, -4, -5})
//	outputs := helper.Duplicates[float64](helper.SliceToChan(expected), 2)
//
//	fmt.Println(helper.ChanToSlice(outputs[0])) // [-10, 20, -4, -5]
//	fmt.Println(helper.ChanToSlice(outputs[1])) // [-10, 20, -4, -5]
func Duplicate[T any](input <-chan T, count int) []<-chan T {
	// TODO(cinar): Find a way to cast as a directional channel array.
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

		for n := range input {
			for _, output := range outputs {
				output <- n
			}
		}
	}()

	return result
}
