// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"io"
	"log"
)

// CloseAndLogError attempts to close the closer and logs any error.
func CloseAndLogError(closer io.Closer, message string) {
	err := closer.Close()
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}
