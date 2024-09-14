// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestDaysBetween(t *testing.T) {
	from := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC)

	actual := helper.DaysBetween(from, from)
	expected := 0

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}

	actual = helper.DaysBetween(from, to)
	expected = 14

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}
