// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"encoding/json"
	"io"
	"log/slog"
)

// JSONToChan reads values from the specified reader in JSON format into a channel of values.
func JSONToChan[T any](r io.Reader) <-chan T {
	return JSONToChanWithLogger[T](r, slog.Default())
}

// JSONToChanWithLogger reads values from the specified reader in JSON format into a channel of values.
func JSONToChanWithLogger[T any](r io.Reader, logger *slog.Logger) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		decoder := json.NewDecoder(r)

		token, err := decoder.Token()
		if err != nil {
			logger.Error("Unable to read token.", "error", err)
			return
		}

		if token != json.Delim('[') {
			logger.Error("Expecting start of array.", "token", token)
			return
		}

		for decoder.More() {
			var value T

			err = decoder.Decode(&value)
			if err != nil {
				logger.Error("Unable to decode value.", "error", err)
				return
			}

			c <- value
		}

		token, err = decoder.Token()
		if err != nil {
			logger.Error("Unable to read token.", "error", err)
			return
		}

		if token != json.Delim(']') {
			logger.Error("Expecting end of array.", "token", token)
			return
		}
	}()

	return c
}
