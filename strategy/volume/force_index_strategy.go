// Copyright (c) 2021-2026 Onur Cinar.
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

// ForceIndexStrategy represents the configuration parameters for calculating the Force Index strategy.
// It recommends a Buy action when it crosses above zero, and a Sell action when it crosses below zero.
type ForceIndexStrategy struct {
	// ForceIndex is the Force Index instance.
	ForceIndex *volume.Fi[float64]
}

// NewForceIndexStrategy function initializes a new Force Index strategy instance with the default parameters.
func NewForceIndexStrategy() *ForceIndexStrategy {
	return NewForceIndexStrategyWith(
		volume.DefaultFiPeriod,
	)
}

// NewForceIndexStrategyWith function initializes a new Force Index strategy instance with the given parameters.
func NewForceIndexStrategyWith(period int) *ForceIndexStrategy {
	return &ForceIndexStrategy{
		ForceIndex: volume.NewFiWithPeriod[float64](period),
	}
}

// Name returns the name of the strategy.
func (f *ForceIndexStrategy) Name() string {
	return fmt.Sprintf("Force Index Strategy (%d)", f.ForceIndex.IdlePeriod()+1)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (f *ForceIndexStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	closings := asset.SnapshotsAsClosings(snapshotsSplice[0])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[1])

	fis := f.ForceIndex.Compute(closings, volumes)

	actions := helper.Map(fis, func(fi float64) strategy.Action {
		if fi > 0 {
			return strategy.Buy
		}

		if fi < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Force Index starts only after a full period.
	actions = helper.Shift(actions, f.ForceIndex.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (f *ForceIndexStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> force index
	// snapshots[2] -> volumes
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	dates := helper.Skip(asset.SnapshotsAsDates(snapshots[0]), f.ForceIndex.IdlePeriod())

	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots[1]),
		2,
	)
	volumes := asset.SnapshotsAsVolumes(snapshots[2])

	fis := f.ForceIndex.Compute(closingsSplice[0], volumes)

	closingsSplice[1] = helper.Skip(closingsSplice[1], f.ForceIndex.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(f, snapshots[3])
	actions = helper.Skip(actions, f.ForceIndex.IdlePeriod())
	outcomes = helper.Skip(outcomes, f.ForceIndex.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(f.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Force Index", fis), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
