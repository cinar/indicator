// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// StopLossStrategy prevents a loss by recommending a sell action when the assets drops below the given threshold.
type StopLossStrategy struct {
	strategy.Strategy

	// InnertStrategy is the inner strategy.
	InnertStrategy strategy.Strategy

	// Percentage is the loss threshold in percentage.
	Percentage float64
}

// NewStopLossStrategy function initializes a new stop loss strategy instance.
func NewStopLossStrategy(innerStrategy strategy.Strategy, percentage float64) *StopLossStrategy {
	return &StopLossStrategy{
		InnertStrategy: innerStrategy,
		Percentage:     percentage,
	}
}

// Name returns the name of the strategy.
func (s *StopLossStrategy) Name() string {
	return fmt.Sprintf("Stop Loss Strategy (%s)", s.InnertStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *StopLossStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	innerActions := s.InnertStrategy.Compute(snapshotsSplice[0])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[1])
	stopLossAt := 0.0

	return helper.Operate(innerActions, closings, func(action strategy.Action, closing float64) strategy.Action {
		// If action is Buy and the asset is not yet bought, buy it as recommended.
		if action == strategy.Buy && stopLossAt == 0.0 {
			stopLossAt = closing * (1 - s.Percentage)
			return strategy.Buy
		}

		// If asset is bought and action is sell or closing is less than or equal to stop loss at, recommend sell.
		if stopLossAt != 0 && (action == strategy.Sell || closing <= stopLossAt) {
			stopLossAt = 0.0
			return strategy.Sell
		}

		return strategy.Hold
	})
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (s *StopLossStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
