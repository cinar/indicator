// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"time"

	"github.com/cinar/indicator/v2/helper"
)

// Snapshot captures a single observation of an asset's price
// at a specific moment.
type Snapshot struct {
	// Date represents the specific timestamp.
	Date time.Time

	// Open represents the opening price for the
	// snapshot period.
	Open float64

	// High represents the highest price reached
	// during the snapshot period.
	High float64

	// Low represents the lowest price reached
	// during the snapshot period.
	Low float64

	// Close represents the closing price for the
	// snapshot period.
	Close float64

	// Volume represents the total trading activity for
	// the asset during the snapshot period.
	Volume float64
}

// SnapshotsAsDates extracts the date field from each snapshot in the provided
// channel and returns a new channel containing only those date values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsDates(snapshots <-chan *Snapshot) <-chan time.Time {
	return helper.Map(snapshots, func(snapshot *Snapshot) time.Time {
		return snapshot.Date
	})
}

// SnapshotsAsOpenings extracts the open field from each snapshot in the provided
// channel and returns a new channel containing only those open values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsOpenings(snapshots <-chan *Snapshot) <-chan float64 {
	return helper.Map(snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Open
	})
}

// SnapshotsAsHighs extracts the high field from each snapshot in the provided
// channel and returns a new channel containing only those high values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsHighs(snapshots <-chan *Snapshot) <-chan float64 {
	return helper.Map(snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.High
	})
}

// SnapshotsAsLows extracts the low field from each snapshot in the provided
// channel and returns a new channel containing only those low values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsLows(snapshots <-chan *Snapshot) <-chan float64 {
	return helper.Map(snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Low
	})
}

// SnapshotsAsClosings extracts the close field from each snapshot in the provided
// channel and returns a new channel containing only those close values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsClosings(snapshots <-chan *Snapshot) <-chan float64 {
	return helper.Map(snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Close
	})
}

// SnapshotsAsVolumes extracts the volume field from each snapshot in the provided
// channel and returns a new channel containing only those volume values.The
// original snapshots channel can no longer be directly used afterward.
func SnapshotsAsVolumes(snapshots <-chan *Snapshot) <-chan float64 {
	return helper.Map(snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Volume
	})
}
