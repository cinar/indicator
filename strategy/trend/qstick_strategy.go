// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// QstickStrategy represents the configuration parameters for calculating the
// Qstick strategy. Qstick is a momentum indicator used to identify
// an asset's trend by looking at the SMA of the difference between
// its closing and opening.
//
// A Qstick above zero indicates increasing buying pressure, while
// a Qstick below zero indicates increasing selling pressure.
type QstickStrategy struct {
	strategy.Strategy

	// Qstick represents the configuration parameters for calculating the Qstick.
	Qstick *momentum.Qstick[float64]
}

// NewQstickStrategy function initializes a new Qstick strategy instance.
func NewQstickStrategy() *QstickStrategy {
	return &QstickStrategy{
		Qstick: momentum.NewQstick[float64](),
	}
}

// Name returns the name of the strategy.
func (*QstickStrategy) Name() string {
	return "Qstick Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (q *QstickStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshots := helper.Duplicate(c, 2)
	openings := asset.SnapshotsAsOpenings(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	qstick := q.Qstick.Compute(openings, closings)
	qstick = helper.Buffered(qstick, 2)

	qsticks := helper.Duplicate(qstick, 2)
	qsticks[1] = helper.Skip(qsticks[1], 1)

	actions := helper.Operate(qsticks[0], qsticks[1], func(b, c float64) strategy.Action {
		// A Qstick above zero indicates increasing buying pressure.
		if c >= 0 && b < 0 {
			return strategy.Buy
		}

		// A Qstick below zero indicates increasing selling pressure.
		if c <= 0 && b > 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Qstick starts only after a full period.
	actions = helper.Shift(actions, q.Qstick.Sma.Period, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (q *QstickStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> openings[1] -> openings
	//                 openings[0] |
	// snapshots[2] -> closings[0] |> qstick
	//                 closings[1] -> closings
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	dates := asset.SnapshotsAsDates(snapshots[0])
	openings := helper.Duplicate(asset.SnapshotsAsOpenings(snapshots[1]), 2)
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 2)

	qstick := q.Qstick.Compute(openings[0], closings[0])
	qstick = helper.Shift(qstick, q.Qstick.Sma.Period-1, 0)

	actions, outcomes := strategy.ComputeWithOutcome(q, snapshots[3])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(q.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Open", openings[1]))
	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("Qstick", qstick), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
