// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"os"
	"testing"
)

// Remove removes the file with the given name.
func Remove(t *testing.T, name string) {
	t.Helper()

	err := os.Remove(name)
	if err != nil {
		t.Errorf("Error removing file: %v", err)
	}
}

// RemoveAll removes the files with the given path.
func RemoveAll(t *testing.T, path string) {
	t.Helper()

	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("Error removing files: %v", err)
	}
}
