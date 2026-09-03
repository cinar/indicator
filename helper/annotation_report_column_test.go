// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestAnnotationReportColumnValue(t *testing.T) {
	values := helper.SliceToChan([]string{"buy", ""})
	column := helper.NewAnnotationReportColumn(values)

	value, err := column.Value()
	if err != nil {
		t.Fatal(err)
	}

	if value != `"buy"` {
		t.Fatalf("value is %s, expected \"buy\"", value)
	}

	value, err = column.Value()
	if err != nil {
		t.Fatal(err)
	}

	if value != "null" {
		t.Fatalf("value is %s, expected null", value)
	}
}

func TestAnnotationReportColumnValueExhausted(t *testing.T) {
	values := helper.SliceToChan([]string{"buy"})
	column := helper.NewAnnotationReportColumn(values)

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
