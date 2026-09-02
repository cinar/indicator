// Copyright (c) 2021-2026 The Indicator Authors.
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

func TestDaysBetweenAcrossDaylightSavingTransition(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	// Spring-forward transition on 2024-03-10 makes this span 71
	// elapsed hours, not the 72 that 3 calendar days implies.
	from := time.Date(2024, 3, 8, 0, 0, 0, 0, loc)
	to := time.Date(2024, 3, 11, 0, 0, 0, 0, loc)

	actual := helper.DaysBetween(from, to)
	expected := 3

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestDaysBetweenNonAlignedTimeOfDay(t *testing.T) {
	from := time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 1, 0, 0, 0, time.UTC)

	actual := helper.DaysBetween(from, to)
	expected := 1

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}
