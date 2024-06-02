// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"encoding/json"
	"io"
)

// ChanToJSON converts a channel of values into JSON format and writes it to the specified writer.
//
// Example:
//
//	input := helper.SliceToChan([]int{2, 4, 6, 8})
//
//	var buffer bytes.Buffer
//	err := helper.ChanToJSON(input, &buffer)
//
//	fmt.Println(buffer.String())
//	// Output: [2,4,6,8,9]
func ChanToJSON[T any](c <-chan T, w io.Writer) error {
	first := true

	_, err := w.Write([]byte{'['})
	if err != nil {
		return err
	}

	for n := range c {
		if !first {
			_, err = w.Write([]byte{','})
			if err != nil {
				return err
			}
		} else {
			first = false
		}

		encoded, err := json.Marshal(n)
		if err != nil {
			return err
		}

		_, err = w.Write(encoded)
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte{']'})

	return err
}
