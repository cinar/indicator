// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"context"
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

// SnapshotsAsDatesWithContext extracts the date field from each snapshot in the provided
// channel and returns a new channel containing only those date values, supporting context cancellation.
func SnapshotsAsDatesWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan time.Time {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) time.Time {
		return snapshot.Date
	})
}

// SnapshotsAsDates extracts the date field from each snapshot in the provided
// channel and returns a new channel containing only those date values.
//
// Deprecated: Use SnapshotsAsDatesWithContext instead.
func SnapshotsAsDates(snapshots <-chan *Snapshot) <-chan time.Time {
	return SnapshotsAsDatesWithContext(context.Background(), snapshots)
}

// SnapshotsAsOpeningsWithContext extracts the open field from each snapshot in the provided
// channel and returns a new channel containing only those open values, supporting context cancellation.
func SnapshotsAsOpeningsWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan float64 {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Open
	})
}

// SnapshotsAsOpenings extracts the open field from each snapshot in the provided
// channel and returns a new channel containing only those open values.
//
// Deprecated: Use SnapshotsAsOpeningsWithContext instead.
func SnapshotsAsOpenings(snapshots <-chan *Snapshot) <-chan float64 {
	return SnapshotsAsOpeningsWithContext(context.Background(), snapshots)
}

// SnapshotsAsHighsWithContext extracts the high field from each snapshot in the provided
// channel and returns a new channel containing only those high values, supporting context cancellation.
func SnapshotsAsHighsWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan float64 {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.High
	})
}

// SnapshotsAsHighs extracts the high field from each snapshot in the provided
// channel and returns a new channel containing only those high values.
//
// Deprecated: Use SnapshotsAsHighsWithContext instead.
func SnapshotsAsHighs(snapshots <-chan *Snapshot) <-chan float64 {
	return SnapshotsAsHighsWithContext(context.Background(), snapshots)
}

// SnapshotsAsLowsWithContext extracts the low field from each snapshot in the provided
// channel and returns a new channel containing only those low values, supporting context cancellation.
func SnapshotsAsLowsWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan float64 {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Low
	})
}

// SnapshotsAsLows extracts the low field from each snapshot in the provided
// channel and returns a new channel containing only those low values.
//
// Deprecated: Use SnapshotsAsLowsWithContext instead.
func SnapshotsAsLows(snapshots <-chan *Snapshot) <-chan float64 {
	return SnapshotsAsLowsWithContext(context.Background(), snapshots)
}

// SnapshotsAsClosingsWithContext extracts the close field from each snapshot in the provided
// channel and returns a new channel containing only those close values, supporting context cancellation.
func SnapshotsAsClosingsWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan float64 {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Close
	})
}

// SnapshotsAsClosings extracts the close field from each snapshot in the provided
// channel and returns a new channel containing only those close values.
//
// Deprecated: Use SnapshotsAsClosingsWithContext instead.
func SnapshotsAsClosings(snapshots <-chan *Snapshot) <-chan float64 {
	return SnapshotsAsClosingsWithContext(context.Background(), snapshots)
}

// SnapshotsAsVolumesWithContext extracts the volume field from each snapshot in the provided
// channel and returns a new channel containing only those volume values, supporting context cancellation.
func SnapshotsAsVolumesWithContext(ctx context.Context, snapshots <-chan *Snapshot) <-chan float64 {
	return helper.MapWithContext(ctx, snapshots, func(snapshot *Snapshot) float64 {
		return snapshot.Volume
	})
}

// SnapshotsAsVolumes extracts the volume field from each snapshot in the provided
// channel and returns a new channel containing only those volume values.
//
// Deprecated: Use SnapshotsAsVolumesWithContext instead.
func SnapshotsAsVolumes(snapshots <-chan *Snapshot) <-chan float64 {
	return SnapshotsAsVolumesWithContext(context.Background(), snapshots)
}
