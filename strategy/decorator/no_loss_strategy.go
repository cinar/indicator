// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// NoLossStrategy prevents selling an asset at a loss. It modifies the recommendations of another strategy to ensure
// that the asset is only sold if its value is at or above the original purchase price.
type NoLossStrategy struct {
	// InnertStrategy is the inner strategy.
	InnertStrategy strategy.Strategy
}

// NewNoLossStrategy function initializes a new no loss strategy instance.
func NewNoLossStrategy(innerStrategy strategy.Strategy) *NoLossStrategy {
	return &NoLossStrategy{
		InnertStrategy: innerStrategy,
	}
}

// Name returns the name of the strategy.
func (n *NoLossStrategy) Name() string {
	return fmt.Sprintf("No Loss Strategy (%s)", n.InnertStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (n *NoLossStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	innerActions := n.InnertStrategy.Compute(snapshotsSplice[0])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[1])
	boughtAt := 0.0

	return helper.Operate(innerActions, closings, func(action strategy.Action, closing float64) strategy.Action {
		// If action is Buy and the asset is not yet bought, buy it as recommended.
		if action == strategy.Buy && boughtAt == 0.0 {
			boughtAt = closing
			return strategy.Buy
		}

		// If the action is sell and the asset was bought at a lower amount, sell it as recommended.
		if action == strategy.Sell && boughtAt != 0.0 && boughtAt < closing {
			boughtAt = 0.0
			return strategy.Sell
		}

		return strategy.Hold
	})
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (n *NoLossStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(n, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(n.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
