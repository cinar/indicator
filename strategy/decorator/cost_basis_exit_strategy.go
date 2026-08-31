// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// CostBasisExitStrategy filters sell recommendations to ensure an asset is only sold if its value is at or above
// the original purchase price (cost basis).
//
// Risk Warning: This algorithmic strategy is a programmatic threshold filter and does NOT guarantee profit or eliminate
// market risk. Holding assets during prolonged or permanent market downturns can result in severe or total loss of capital.
type CostBasisExitStrategy struct {
	// InnerStrategy is the inner strategy.
	InnerStrategy strategy.Strategy
}

// NewCostBasisExitStrategy function initializes a new cost basis exit strategy instance.
func NewCostBasisExitStrategy(innerStrategy strategy.Strategy) *CostBasisExitStrategy {
	return &CostBasisExitStrategy{
		InnerStrategy: innerStrategy,
	}
}

// Name returns the name of the strategy.
func (c *CostBasisExitStrategy) Name() string {
	return fmt.Sprintf("Cost Basis Exit Strategy (%s)", c.InnerStrategy.Name())
}

// ComputeWithContext processes the provided asset snapshots and generates a stream of actionable recommendations.
func (c *CostBasisExitStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 2)

	innerActions := strategy.ComputeStrategyWithContext(ctx, c.InnerStrategy, snapshotsSplice[0])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[1])
	boughtAt := 0.0

	return helper.OperateWithContext(ctx, innerActions, closings, func(action strategy.Action, closing float64) strategy.Action {
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
func (c *CostBasisExitStrategy) Report(s <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(s, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(c, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(c.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (c *CostBasisExitStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return c.ComputeWithContext(context.Background(), snapshots)
}

// NoLossStrategy is an alias for CostBasisExitStrategy for backwards compatibility.
//
// Deprecated: Use CostBasisExitStrategy instead.
type NoLossStrategy = CostBasisExitStrategy

// NewNoLossStrategy is an alias for NewCostBasisExitStrategy for backwards compatibility.
//
// Deprecated: Use NewCostBasisExitStrategy instead.
func NewNoLossStrategy(innerStrategy strategy.Strategy) *CostBasisExitStrategy {
	return NewCostBasisExitStrategy(innerStrategy)
}
