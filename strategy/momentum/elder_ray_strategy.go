// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// ElderRayStrategy represents the configuration parameters for calculating the Elder Ray strategy.
// A positive Bull Power suggests a Buy signal, while a negative Bear Power suggests a Sell signal.
type ElderRayStrategy struct {
	// ElderRay represents the configuration parameters for calculating the Elder-Ray Index.
	ElderRay *momentum.ElderRay[float64]
}

// NewElderRayStrategy function initializes a new Elder Ray strategy instance with the default parameters.
func NewElderRayStrategy() *ElderRayStrategy {
	return &ElderRayStrategy{
		ElderRay: momentum.NewElderRay[float64](),
	}
}

// Name returns the name of the strategy.
func (*ElderRayStrategy) Name() string {
	return "Elder Ray Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (e *ElderRayStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])

	bullPower, bearPower := e.ElderRay.Compute(highs, lows, closings)

	actions := helper.Operate(bullPower, bearPower, func(bull, bear float64) strategy.Action {
		if bull > 0 {
			return strategy.Buy
		}

		if bear < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Elder Ray starts only after the idle period.
	actions = helper.Shift(actions, e.ElderRay.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (e *ElderRayStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs   -|
	// snapshots[2] -> lows    -+-> ElderRay.Compute -> bullPower, bearPower
	// snapshots[3] -> closings-|
	// snapshots[4] -> closings -> close
	// snapshots[5] -> actions  -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := asset.SnapshotsAsClosings(snapshots[3])
	closings2 := asset.SnapshotsAsClosings(snapshots[4])

	bullPower, bearPower := e.ElderRay.Compute(highs, lows, closings)
	bullPower = helper.Shift(bullPower, e.ElderRay.IdlePeriod(), 0)
	bearPower = helper.Shift(bearPower, e.ElderRay.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(e, snapshots[5])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(e.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings2))
	report.AddColumn(helper.NewNumericReportColumn("Bull Power", bullPower), 1)
	report.AddColumn(helper.NewNumericReportColumn("Bear Power", bearPower), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
