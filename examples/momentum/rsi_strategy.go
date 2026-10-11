// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultRsiStrategyBuyAt defines the default RSI level at which a Buy action is generated.
	DefaultRsiStrategyBuyAt = 30

	// DefaultRsiStrategySellAt defines the default RSI level at which a Sell action is generated.
	DefaultRsiStrategySellAt = 70
)

// RsiStrategy demonstrates how to compose the Relative Strength Index (RSI)
// into an illustrative threshold-based overbought/oversold strategy.
type RsiStrategy struct {
	// Rsi represents the configuration parameters for calculating the Relative Strength Index (RSI).
	Rsi *momentum.Rsi[float64]

	// BuyAt defines the RSI level at which a Buy action is generated.
	BuyAt float64

	// SellAt defines the RSI level at which a Sell action is generated.
	SellAt float64
}

// NewRsiStrategy initializes an example RsiStrategy instance with default parameters.
func NewRsiStrategy() *RsiStrategy {
	return NewRsiStrategyWith(
		DefaultRsiStrategyBuyAt,
		DefaultRsiStrategySellAt,
	)
}

// NewRsiStrategyWith initializes an example RsiStrategyWith instance with default parameters.
func NewRsiStrategyWith(buyAt, sellAt float64) *RsiStrategy {
	return &RsiStrategy{
		Rsi:    momentum.NewRsi[float64](),
		BuyAt:  buyAt,
		SellAt: sellAt,
	}
}

// Name returns the name of the example strategy.
func (r *RsiStrategy) Name() string {
	return fmt.Sprintf("RSI Strategy %.0f-%.0f", r.BuyAt, r.SellAt)
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (r *RsiStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots)

	rsi := r.Rsi.ComputeWithContext(ctx, closings)

	actions := helper.MapWithContext(ctx, rsi, func(value float64) strategy.Action {
		if value <= r.BuyAt {
			return strategy.Buy
		}

		if value >= r.SellAt {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// RSI starts only after the idle period.
	actions = helper.ShiftWithContext(ctx, actions, r.Rsi.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (r *RsiStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute     -> actions -> annotations
	// snapshots[2] -> closings[0] -> close
	//              -> closings[1] -> Rsi.Compute -> rsi
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 2)
	rsi := helper.Shift(r.Rsi.Compute(closings[1]), r.Rsi.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(r, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(r.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("RSI", rsi), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (r *RsiStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return r.ComputeWithContext(context.Background(), snapshots)
}
