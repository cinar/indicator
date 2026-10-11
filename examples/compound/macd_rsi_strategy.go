// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package compound

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/examples/momentum"
	"github.com/cinar/indicator/v2/examples/trend"
)

const (
	// DefaultMacdRsiStrategyBuyAt defines the default RSI level at which a Buy action is generated.
	DefaultMacdRsiStrategyBuyAt = 30

	// DefaultMacdRsiStrategySellAt defines the default RSI level at which a Sell action is generated.
	DefaultMacdRsiStrategySellAt = 70
)

// MacdRsiStrategy demonstrates how to compose MACD and RSI example strategies
// into an illustrative combined strategy.
type MacdRsiStrategy struct {
	// MacdStrategy is the MACD strategy instance.
	MacdStrategy *trend.MacdStrategy

	// RsiStrategy is the RSI strategy instance.
	RsiStrategy *momentum.RsiStrategy
}

// NewMacdRsiStrategy initializes an example MacdRsiStrategy instance with default parameters.
func NewMacdRsiStrategy() *MacdRsiStrategy {
	return NewMacdRsiStrategyWith(
		DefaultMacdRsiStrategyBuyAt,
		DefaultMacdRsiStrategySellAt,
	)
}

// NewMacdRsiStrategyWith initializes an example MacdRsiStrategyWith instance with default parameters.
func NewMacdRsiStrategyWith(buyAt, sellAt float64) *MacdRsiStrategy {
	return &MacdRsiStrategy{
		MacdStrategy: trend.NewMacdStrategy(),
		RsiStrategy:  momentum.NewRsiStrategyWith(buyAt, sellAt),
	}
}

// Name returns the name of the example strategy.
func (m *MacdRsiStrategy) Name() string {
	return fmt.Sprintf("MACD-RSI Strategy (%.0f, %.0f)",
		m.RsiStrategy.BuyAt,
		m.RsiStrategy.SellAt,
	)
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (m *MacdRsiStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 2)

	macdActions := strategy.DenormalizeActions(
		m.MacdStrategy.ComputeWithContext(ctx, snapshotsSplice[0]),
	)

	rsiActions := strategy.DenormalizeActions(
		m.RsiStrategy.ComputeWithContext(ctx, snapshotsSplice[1]),
	)

	actions := helper.OperateWithContext(ctx, macdActions, rsiActions, func(macdAction, rsiAction strategy.Action) strategy.Action {
		if macdAction == rsiAction {
			return macdAction
		}

		return strategy.Hold
	})

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (m *MacdRsiStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> macds, signals
	//                 closings[2] -> rsi
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 3)

	macds, signals := m.MacdStrategy.Macd.Compute(closings[0])
	macds = helper.Shift(macds, m.MacdStrategy.Macd.IdlePeriod(), 0)
	signals = helper.Shift(signals, m.MacdStrategy.Macd.IdlePeriod(), 0)

	rsi := m.RsiStrategy.Rsi.Compute(closings[2])
	rsi = helper.Shift(rsi, m.RsiStrategy.Rsi.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(m, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(m.Name(), dates)
	report.AddChart()
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("MACD", macds), 1)
	report.AddColumn(helper.NewNumericReportColumn("Signal", signals), 1)

	report.AddColumn(helper.NewNumericReportColumn("RSI", rsi), 2)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1, 2)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 3)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *MacdRsiStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return m.ComputeWithContext(context.Background(), snapshots)
}
