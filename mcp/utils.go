package main

import (
	"time"
)

// toTimeArray converts a slice of Unix timestamps (seconds since epoch) to a slice of time.Time values.
// This is useful for converting timestamp data from JSON/API responses into Go's native time type.
//
// Parameters:
//   - timestamps: A slice of int64 values representing Unix timestamps in seconds
//
// Returns:
//   - A slice of time.Time values corresponding to the input timestamps
//   - The order of elements in the output matches the order of the input timestamps
func toTimeArray(timestamps []int64) []time.Time {
	timeArray := make([]time.Time, len(timestamps))
	for i, ts := range timestamps {
		timeArray[i] = time.Unix(ts, 0)
	}
	return timeArray
}
