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

// CciStrategy represents the configuration parameters for calculating the CCI strategy.
// A CCI value crossing above the 100+ suggests a bullish trend, while crossing below
// the 100- indicates a bearish trend.
type CciStrategy struct {
	// Cci represents the configuration parameters for calculating the CCI.
	Cci *trend.Cci[float64]
}

// NewCciStrategy function initializes a new CCI strategy instance.
func NewCciStrategy() *CciStrategy {
	return &CciStrategy{
		Cci: trend.NewCci[float64](),
	}
}

// Name returns the name of the strategy.
func (*CciStrategy) Name() string {
	return "CCI Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (t *CciStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshots := helper.Duplicate(c, 3)
	highs := asset.SnapshotsAsHighs(snapshots[0])
	lows := asset.SnapshotsAsLows(snapshots[1])
	closings := asset.SnapshotsAsClosings(snapshots[2])

	ccis := t.Cci.Compute(highs, lows, closings)

	actions := helper.Map(ccis, func(cci float64) strategy.Action {
		if cci >= 100 {
			return strategy.Buy
		}

		if cci <= -100 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// CCI starts only after a full period.
	actions = helper.Shift(actions, t.Cci.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (t *CciStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[1] |> ccis
	//                 closings[0] -> closings
	// snapshots[4] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[3]), 2)

	ccis := t.Cci.Compute(highs, lows, closings[1])
	ccis = helper.Shift(ccis, t.Cci.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("CCI", ccis), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
