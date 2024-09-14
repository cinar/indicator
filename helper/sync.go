// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "slices"

// CommonPeriod calculates the smallest period at which all data channels can be synchronized
//
// Example:
//
//	// Synchronize channels with periods 4, 2, and 3.
//	commonPeriod := helper.CommonPeriod(4, 2, 3) // commonPeriod = 4
//
//	// Synchronize the first channel
//	c1 := helper.Sync(commonPeriod, 4, c1)
//
//	// Synchronize the second channel
//	c2 := helper.Sync(commonPeriod, 2, c2)
//
//	// Synchronize the third channel
//	c3 := helper.Sync(commonPeriod, 3, c3)
func CommonPeriod(periods ...int) int {
	return slices.Max(periods)
}

// SyncPeriod adjusts the given channel to match the given common period.
func SyncPeriod[T any](commonPeriod, period int, c <-chan T) <-chan T {
	forwardPeriod := commonPeriod - period

	if forwardPeriod > 0 {
		c = Skip(c, forwardPeriod)
	}

	return c
}
