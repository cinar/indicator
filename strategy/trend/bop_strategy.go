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

// BopStrategy gauges the strength of buying and selling forces using the
// Balance of Power (BoP) indicator. A positive BoP value  suggests an
// upward trend, while a negative value indicates a downward trend. A
// BoP value of zero implies equilibrium between the two forces.
type BopStrategy struct {
	strategy.Strategy

	// Bop represents the configuration parameters for calculating the
	// Balance of Power (BoP).
	Bop *trend.Bop[float64]
}

// NewBopStrategy function initializes a new BoP strategy instance with the default parameters.
func NewBopStrategy() *BopStrategy {
	return &BopStrategy{
		Bop: trend.NewBop[float64](),
	}
}

// Name returns the name of the strategy.
func (*BopStrategy) Name() string {
	return "BoP Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (b *BopStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshots := helper.Duplicate(c, 4)

	openings := asset.SnapshotsAsOpenings(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := asset.SnapshotsAsClosings(snapshots[3])

	bops := b.Bop.Compute(openings, highs, lows, closings)

	return strategy.NormalizeActions(helper.Map(bops, func(bop float64) strategy.Action {
		if bop > 0 {
			return strategy.Buy
		}

		if bop < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	}))
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (b *BopStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> openings    |
	// snapshots[2] -> highs       |
	// snapshots[3] -> lows        |
	// snapshots[4] -> closings[1] |> bop
	//                 closings[0] -> closings
	// snapshots[5] -> actions     -> annotations
	//                 outcomes
	//
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	openings := asset.SnapshotsAsOpenings(snapshots[1])
	highs := asset.SnapshotsAsHighs(snapshots[2])
	lows := asset.SnapshotsAsLows(snapshots[3])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[4]), 2)

	bop := b.Bop.Compute(openings, highs, lows, closings[1])

	actions, outcomes := strategy.ComputeWithOutcome(b, snapshots[5])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(b.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("BoP", bop), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
