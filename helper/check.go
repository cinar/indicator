// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"errors"
	"fmt"
	"reflect"
)

// CheckEquals determines whether the two channels are equal.
func CheckEquals[T comparable](inputs ...<-chan T) error {
	if len(inputs)%2 != 0 {
		return errors.New("not pairs")
	}

	i := 0

	for {
		for j := 0; j < len(inputs); j += 2 {
			actual, actualOk := <-inputs[j]
			expected, expectedOk := <-inputs[j+1]

			if !actualOk || !expectedOk {
				if actualOk != expectedOk {
					return fmt.Errorf("not ended the same actual %v expected %v", actualOk, expectedOk)
				}

				return nil
			}

			if !reflect.DeepEqual(actual, expected) {
				return fmt.Errorf("index %d pair %d actual %v expected %v", i, j/2, actual, expected)
			}
		}

		i++
	}
}
