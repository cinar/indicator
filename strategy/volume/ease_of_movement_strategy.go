// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volume"
)

// EaseOfMovementStrategy represents the configuration parameters for calculating the Ease of Movement strategy.
// Recommends a Buy action when it crosses above 0, and recommends a Sell action when it crosses below 0.
type EaseOfMovementStrategy struct {
	// EaseOfMovement is the Ease of Movement indicator instance.
	EaseOfMovement *volume.Emv[float64]
}

// NewEaseOfMovementStrategy function initializes a new Ease of Movement strategy instance with the
// default parameters.
func NewEaseOfMovementStrategy() *EaseOfMovementStrategy {
	return NewEaseOfMovementStrategyWith(
		volume.DefaultEmvPeriod,
	)
}

// NewEaseOfMovementStrategyWith function initializes a new Ease of Movement strategy instance with the
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

// Compute function processes the provided asset snapshots and generates a stream of actionable recommendations.
func (e *EaseOfMovementStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[2])

	emvs := e.EaseOfMovement.Compute(highs, lows, volumes)

	actions := helper.Map(emvs, func(emv float64) strategy.Action {
		if emv > 0 {
			return strategy.Buy
		}

		if emv < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Ease of Movement starts only after a full period.
	actions = helper.Shift(actions, e.EaseOfMovement.IdlePeriod(), strategy.Hold)

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
