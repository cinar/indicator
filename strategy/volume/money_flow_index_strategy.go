// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volume"
)

const (
	// DefaultMoneyFlowIndexStrategySellAt is the default sell at of 80.
	DefaultMoneyFlowIndexStrategySellAt = 80

	// DefaultMoneyFlowIndexStrategyBuyAt is the default buy at of 20.
	DefaultMoneyFlowIndexStrategyBuyAt = 20
)

// MoneyFlowIndexStrategy represents the configuration parameters for calculating the Money Flow Index strategy.
// Recommends a Sell action when it crosses over 80, and recommends a Buy action when it crosses below 20.
type MoneyFlowIndexStrategy struct {
	// MoneyFlowIndex is the Money Flow Index indicator instance.
	MoneyFlowIndex *volume.Mfi[float64]

	// SellAt is the sell at value.
	SellAt float64

	// BuyAt is the buy at value.
	BuyAt float64
}

// NewMoneyFlowIndexStrategy function initializes a new Money Flow Index strategy instance with the default parameters.
func NewMoneyFlowIndexStrategy() *MoneyFlowIndexStrategy {
	return NewMoneyFlowIndexStrategyWith(
		DefaultMoneyFlowIndexStrategySellAt,
		DefaultMoneyFlowIndexStrategyBuyAt,
	)
}

// NewMoneyFlowIndexStrategyWith function initializes a new Money Flow Index strategy instance with the
// given parameters.
func NewMoneyFlowIndexStrategyWith(sellAt, buyAt float64) *MoneyFlowIndexStrategy {
	return &MoneyFlowIndexStrategy{
		MoneyFlowIndex: volume.NewMfi[float64](),
		SellAt:         sellAt,
		BuyAt:          buyAt,
	}
}

// Name returns the name of the strategy.
func (m *MoneyFlowIndexStrategy) Name() string {
	return fmt.Sprintf("Money Flow Index Strategy (%f)", m.SellAt)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (m *MoneyFlowIndexStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 4)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[3])

	mfis := m.MoneyFlowIndex.Compute(highs, lows, closings, volumes)

	actions := helper.Map(mfis, func(mfi float64) strategy.Action {
		if mfi >= m.SellAt {
			return strategy.Sell
		}

		if mfi <= m.BuyAt {
			return strategy.Buy
		}

		return strategy.Hold
	})

	// Money Flow Index starts only after a full period.
	actions = helper.Shift(actions, m.MoneyFlowIndex.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (m *MoneyFlowIndexStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[0] -> closings
	//                 closings[1] -> superTrend
	// snapshots[4] -> volumes
	// snapshots[5] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 6)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshots[0]),
		m.MoneyFlowIndex.IdlePeriod(),
	)

	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots[3]),
		2,
	)
	volumes := asset.SnapshotsAsVolumes(snapshots[4])

	mfis := m.MoneyFlowIndex.Compute(highs, lows, closingsSplice[0], volumes)
	closingsSplice[1] = helper.Skip(closingsSplice[1], m.MoneyFlowIndex.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(m, snapshots[5])
	actions = helper.Skip(actions, m.MoneyFlowIndex.IdlePeriod())
	outcomes = helper.Skip(outcomes, m.MoneyFlowIndex.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(m.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Money Flow Index", mfis), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
