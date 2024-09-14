// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"math"
	"time"
)

// DaysBetween calculates the days between the given two times.
func DaysBetween(from, to time.Time) int {
	return int(math.Floor(to.Sub(from).Hours() / 24))
}
