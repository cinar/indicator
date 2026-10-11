// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// QstickStrategy demonstrates how to compose the Qstick momentum indicator
// into an illustrative zero-line crossover strategy.
type QstickStrategy struct {
	// Qstick represents the configuration parameters for calculating the Qstick.
	Qstick *momentum.Qstick[float64]
}

// NewQstickStrategy initializes an example QstickStrategy instance with default parameters.
func NewQstickStrategy() *QstickStrategy {
	return &QstickStrategy{
		Qstick: momentum.NewQstick[float64](),
	}
}

// Name returns the name of the example strategy.
func (*QstickStrategy) Name() string {
	return "Qstick Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an
// illustrative stream of actions.
func (q *QstickStrategy) ComputeWithContext(ctx context.Context, c <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshots := helper.DuplicateWithContext(ctx, c, 2)
	openings := asset.SnapshotsAsOpeningsWithContext(ctx, snapshots[0])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots[1])

	qstick := q.Qstick.ComputeWithContext(ctx, openings, closings)
	qstick = helper.BufferedWithContext(ctx, qstick, 2)

	qsticks := helper.DuplicateWithContext(ctx, qstick, 2)
	qsticks[1] = helper.SkipWithContext(ctx, qsticks[1], 1)

	actions := helper.OperateWithContext(ctx, qsticks[0], qsticks[1], func(b, c float64) strategy.Action {
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
	actions = helper.ShiftWithContext(ctx, actions, q.Qstick.Sma.Period, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (q *QstickStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	return q.ComputeWithContext(context.Background(), c)
}
