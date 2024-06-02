// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

// TrixStrategy represents the configuration parameters for calculating the TRIX strategy.
// A TRIX value crossing above the zero line suggests a bullish trend, while crossing
// below the zero line indicates a bearish trend.
type TrixStrategy struct {
	strategy.Strategy

	// Trix represents the configuration parameters for calculating the TRIX.
	Trix *trend.Trix[float64]
}

// NewTrixStrategy function initializes a new TRIX strategy instance.
func NewTrixStrategy() *TrixStrategy {
	return &TrixStrategy{
		Trix: trend.NewTrix[float64](),
	}
}

// Name returns the name of the strategy.
func (*TrixStrategy) Name() string {
	return "TRIX Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (t *TrixStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosings(snapshots)

	trixs := t.Trix.Compute(closings)

	actions := strategy.NormalizeActions(helper.Map(trixs, func(trix float64) strategy.Action {
		if trix > 0 {
			return strategy.Buy
		}

		if trix < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	}))

	// TRIX starts only after a full period.
	actions = helper.Shift(actions, t.Trix.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (t *TrixStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> trixs
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 2)

	trixs := t.Trix.Compute(closings[0])
	trixs = helper.Shift(trixs, t.Trix.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("TRIX", trixs), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
