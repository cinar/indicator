// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"io"
	"log/slog"
)

// CloseAndLogError attempts to close the closer and logs any error.
func CloseAndLogError(closer io.Closer, message string) {
	CloseAndLogErrorWithLogger(closer, message, slog.Default())
}

// CloseAndLogErrorWithLogger attempts to close the closer and logs any error to the given logger.
func CloseAndLogErrorWithLogger(closer io.Closer, message string, logger *slog.Logger) {
	err := closer.Close()
	if err != nil {
		logger.Error(message, "error", err)
	}
}
