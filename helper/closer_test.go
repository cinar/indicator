// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"os"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestCloseAndLogErrorWithoutError(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "closer")
	if err != nil {
		t.Fatal(err)
	}

	helper.CloseAndLogError(file, "")
	helper.CloseAndLogError(file, "")
}
