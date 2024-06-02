// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"encoding/json"
	"io"
	"log"
)

// JSONToChan reads values from the specified reader in JSON format into a channel of values.
//
// Example:
func JSONToChan[T any](r io.Reader) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		decoder := json.NewDecoder(r)

		token, err := decoder.Token()
		if err != nil {
			log.Print(err)
			return
		}

		if token != json.Delim('[') {
			log.Printf("expecting start of array got %v", token)
			return
		}

		for decoder.More() {
			var value T

			err = decoder.Decode(&value)
			if err != nil {
				log.Print(err)
				return
			}

			c <- value
		}

		token, err = decoder.Token()
		if err != nil {
			log.Print(err)
			return
		}

		if token != json.Delim(']') {
			log.Printf("expecting end of array got %v", token)
			return
		}
	}()

	return c
}
