// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
)

// JSONToChan wraps JSONToChanWithContext for backwards compatibility.
//
// Deprecated: Use JSONToChanWithContext instead.
func JSONToChan[T any](r io.Reader) <-chan T {
	return JSONToChanWithContext[T](context.Background(), r)
}

// JSONToChanWithContext reads values from the specified reader in JSON format into a channel of values, supporting context cancellation.
func JSONToChanWithContext[T any](ctx context.Context, r io.Reader) <-chan T {
	return JSONToChanWithLoggerWithContext[T](ctx, r, slog.Default())
}

// JSONToChanWithLogger wraps JSONToChanWithLoggerWithContext for backwards compatibility.
//
// Deprecated: Use JSONToChanWithLoggerWithContext instead.
func JSONToChanWithLogger[T any](r io.Reader, logger *slog.Logger) <-chan T {
	return JSONToChanWithLoggerWithContext[T](context.Background(), r, logger)
}

// JSONToChanWithLoggerWithContext reads values from the specified reader in JSON format into a channel of values with logger and context.
func JSONToChanWithLoggerWithContext[T any](ctx context.Context, r io.Reader, logger *slog.Logger) <-chan T {
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
			select {
			case <-ctx.Done():
				return
			default:
			}

			var value T

			err = decoder.Decode(&value)
			if err != nil {
				logger.Error("Unable to decode value.", "error", err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case c <- value:
			}
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
