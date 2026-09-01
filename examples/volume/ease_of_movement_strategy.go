// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volume"
)

// EaseOfMovementStrategy demonstrates how to compose the Ease of Movement (EMV)
// indicator into an illustrative zero-line crossover strategy.
type EaseOfMovementStrategy struct {
	// EaseOfMovement is the Ease of Movement indicator instance.
	EaseOfMovement *volume.Emv[float64]
}

// NewEaseOfMovementStrategy initializes an example EaseOfMovementStrategy instance with default parameters.
// default parameters.
func NewEaseOfMovementStrategy() *EaseOfMovementStrategy {
	return NewEaseOfMovementStrategyWith(
		volume.DefaultEmvPeriod,
	)
}

// NewEaseOfMovementStrategyWith initializes an example EaseOfMovementStrategyWith instance with default parameters.
// given parameters.
func NewEaseOfMovementStrategyWith(period int) *EaseOfMovementStrategy {
	return &EaseOfMovementStrategy{
		EaseOfMovement: volume.NewEmvWithPeriod[float64](period),
	}
}

// Name function returns the name of the strategy.
func (e *EaseOfMovementStrategy) Name() string {
	return fmt.Sprintf("Ease of Movement Strategy (%d)", e.EaseOfMovement.IdlePeriod()+1)
}

// ComputeWithContext function processes the provided asset snapshots and generates a stream of actionable recommendations.
func (e *EaseOfMovementStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 3)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	volumes := asset.SnapshotsAsVolumesWithContext(ctx, snapshotsSplice[2])

	emvs := e.EaseOfMovement.ComputeWithContext(ctx, highs, lows, volumes)

	actions := helper.MapWithContext(ctx, emvs, func(emv float64) strategy.Action {
		if emv > 0 {
			return strategy.Buy
		}

		if emv < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Ease of Movement starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, e.EaseOfMovement.IdlePeriod(), strategy.Hold)

	return actions
}

// Report function processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (e *EaseOfMovementStrategy) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> volumes     -> emv
	// snapshots[4] -> closings
	// snapshots[5] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(snapshots, 6)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		e.EaseOfMovement.IdlePeriod(),
	)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[1])
	lows := asset.SnapshotsAsLows(snapshotsSplice[2])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[3])

	closings := helper.Skip(
		asset.SnapshotsAsClosings(snapshotsSplice[4]),
		e.EaseOfMovement.IdlePeriod(),
	)

	emvs := e.EaseOfMovement.Compute(highs, lows, volumes)

	actions, outcomes := strategy.ComputeWithOutcome(e, snapshotsSplice[5])
	actions = helper.Skip(actions, e.EaseOfMovement.IdlePeriod())
	outcomes = helper.Skip(outcomes, e.EaseOfMovement.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(e.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("Ease of Movement", emvs), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (e *EaseOfMovementStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return e.ComputeWithContext(context.Background(), snapshots)
}
