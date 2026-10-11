// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestNumericReportColumnValue(t *testing.T) {
	values := helper.SliceToChan([]float64{1, 2})
	column := helper.NewNumericReportColumn("Close", values)

	value, err := column.Value()
	if err != nil {
		t.Fatal(err)
	}

	if value != "1" {
		t.Fatalf("value is %s, expected 1", value)
	}
}

func TestNumericReportColumnValueExhausted(t *testing.T) {
	values := helper.SliceToChan([]float64{1})
	column := helper.NewNumericReportColumn("Close", values)

	// Consume the only available value.
	if _, err := column.Value(); err != nil {
		t.Fatal(err)
	}

	// The backing channel is now exhausted and closed.
	_, err := column.Value()
	if err == nil {
		t.Fatal("expected error")
	}
}
