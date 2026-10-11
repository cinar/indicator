// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volatility"
)

// SuperTrendStrategy demonstrates how to compose the Super Trend indicator
// into an illustrative trend-following strategy.
type SuperTrendStrategy struct {
	// SuperTrend represents the configuration parameters for calculating the Super Trend.
	SuperTrend *volatility.SuperTrend[float64]
}

// NewSuperTrendStrategy initializes an example SuperTrendStrategy instance with default parameters.
func NewSuperTrendStrategy() *SuperTrendStrategy {
	return NewSuperTrendStrategyWith(volatility.NewSuperTrend[float64]())
}

// NewSuperTrendStrategyWith initializes an example SuperTrendStrategyWith instance with default parameters.
func NewSuperTrendStrategyWith(superTrend *volatility.SuperTrend[float64]) *SuperTrendStrategy {
	return &SuperTrendStrategy{
		SuperTrend: superTrend,
	}
}

// Name returns the name of the example strategy.
func (s *SuperTrendStrategy) Name() string {
	return fmt.Sprintf("Super Trend Strategy (%s, %.1f)", s.SuperTrend.Atr.Ma.String(), s.SuperTrend.Multiplier)
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (s *SuperTrendStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 3)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	closingsSplice := helper.DuplicateWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[2]),
		2,
	)

	superTrends := s.SuperTrend.ComputeWithContext(ctx, highs, lows, closingsSplice[0])

	closingsSplice[1] = helper.SkipWithContext(ctx, closingsSplice[1], s.SuperTrend.IdlePeriod())

	actions := helper.OperateWithContext(ctx, superTrends, closingsSplice[1], func(superTrend, closing float64) strategy.Action {
		if superTrend < closing {
			return strategy.Buy
		}

		if superTrend > closing {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Super Trend starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, s.SuperTrend.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (s *SuperTrendStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[0] -> closings
	//                 closings[1] -> superTrend
	// snapshots[4] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots[3]),
		2,
	)

	superTrends := s.SuperTrend.Compute(highs, lows, closingsSplice[0])
	superTrends = helper.Shift(superTrends, s.SuperTrend.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Super Trend", superTrends))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *SuperTrendStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return s.ComputeWithContext(context.Background(), snapshots)
}
