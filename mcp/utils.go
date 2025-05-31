package main

import (
	"time"
)

func toTimeArray(timestamps []int64) []time.Time {
	timeArray := make([]time.Time, len(timestamps))
	for i, ts := range timestamps {
		timeArray[i] = time.Unix(ts, 0)
	}
	return timeArray
}
