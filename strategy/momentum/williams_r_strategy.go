// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultWilliamsRStrategyBuyAt defines the default Williams R level at which a Buy action is generated.
	DefaultWilliamsRStrategyBuyAt = -80.0

	// DefaultWilliamsRStrategySellAt defines the default Williams R level at which a Sell action is generated.
	DefaultWilliamsRStrategySellAt = -20.0
)

// WilliamsRStrategy represents the configuration parameters for calculating the Williams R strategy.
type WilliamsRStrategy struct {
	// WilliamsR represents the configuration parameters for calculating the Williams %R.
	WilliamsR *momentum.WilliamsR[float64]

	// BuyAt defines the Williams R level at which a Buy action is generated.
	BuyAt float64

	// SellAt defines the Williams R level at which a Sell action is generated.
	SellAt float64
}

// NewWilliamsRStrategy function initializes a new Williams R strategy instance with the default parameters.
func NewWilliamsRStrategy() *WilliamsRStrategy {
	return NewWilliamsRStrategyWith(
		DefaultWilliamsRStrategyBuyAt,
		DefaultWilliamsRStrategySellAt,
	)
}

// NewWilliamsRStrategyWith function initializes a new Williams R strategy instance with the given parameters.
func NewWilliamsRStrategyWith(buyAt, sellAt float64) *WilliamsRStrategy {
	return &WilliamsRStrategy{
		WilliamsR: momentum.NewWilliamsR[float64](),
		BuyAt:     buyAt,
		SellAt:    sellAt,
	}
}

// Name returns the name of the strategy.
func (r *WilliamsRStrategy) Name() string {
	return fmt.Sprintf("Williams R Strategy (%.0f,%.0f)", r.BuyAt, r.SellAt)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (r *WilliamsRStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])

	wr := r.WilliamsR.Compute(highs, lows, closings)

	actions := helper.Map(wr, func(value float64) strategy.Action {
		if value <= r.BuyAt {
			return strategy.Buy
		}

		if value >= r.SellAt {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Williams R starts only after the idle period.
	actions = helper.Shift(actions, r.WilliamsR.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (r *WilliamsRStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute          -> actions -> annotations
	// snapshots[2] -> closings         -> close
	// snapshots[3] -> highs   -|
	// snapshots[4] -> lows    -+-> WilliamsR.Compute -> wr
	// snapshots[5] -> closings-|
	//
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[2])
	highs := asset.SnapshotsAsHighs(snapshots[3])
	lows := asset.SnapshotsAsLows(snapshots[4])
	closingsForWR := asset.SnapshotsAsClosings(snapshots[5])

	wr := helper.Shift(r.WilliamsR.Compute(highs, lows, closingsForWR), r.WilliamsR.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(r, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(r.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("Williams R", wr), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
