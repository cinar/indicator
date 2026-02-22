// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volatility"
	"github.com/cinar/indicator/v2/volume"
)

const (
	// DefaultPercentBandMFIStrategyPercentBBuyAt is the default buy for %B at of 0.8.
	DefaultPercentBandMFIStrategyPercentBBuyAt = 0.8

	// DefaultPercentBandMFIStrategyPercentBSellAt is the default sell for %B at of 0.2.
	DefaultPercentBandMFIStrategyPercentBSellAt = 0.2

	// DefaultPercentBandMFIStrategyMfiBuyAt is the default buy for MFI at of 80.
	DefaultPercentBandMFIStrategyMfiBuyAt = 80

	// DefaultPercentBandMFIStrategyMfiSellAt is the default sell for MFI at of 20.
	DefaultPercentBandMFIStrategyMfiSellAt = 20
)

// PercentBandMFIStrategy represents the configuration parameters for calculating the %B combined with MFI strategy.
// Recommends a Buy action when %B is above 0.8 and MFI is above 80, and recommends a Sell action when %B is below 0.2
// and MFI is below 20.
type PercentBandMFIStrategy struct {
	// MoneyFlowIndex is the Money Flow Index indicator instance.
	MoneyFlowIndex *volume.Mfi[float64]

	// PercentB is the %B indicator instance.
	PercentB *volatility.PercentB[float64]

	// SellPercentBAt is the sell at value of %B.
	SellPercentBAt float64

	// BuyPercentBAt is the buy at value of %B.
	BuyPercentBAt float64

	// SellMfiAt is the sell at value of MFI.
	SellMfiAt float64

	// BuyMfiAt is the buy at value of MFI.
	BuyMfiAt float64
}

// NewPercentBandMFIStrategy function initializes a new PercentBandMFI strategy instance with the default parameters.
func NewPercentBandMFIStrategy() *PercentBandMFIStrategy {
	return NewPercentBandMFIStrategyWith(
		DefaultPercentBandMFIStrategyPercentBBuyAt,
		DefaultPercentBandMFIStrategyPercentBSellAt,
		DefaultPercentBandMFIStrategyMfiBuyAt,
		DefaultPercentBandMFIStrategyMfiSellAt,
	)
}

// NewPercentBandMFIStrategyWith function initializes a new PercentBandMFI strategy instance with the
// given parameters.
func NewPercentBandMFIStrategyWith(sellPercentBAt, buyPercentBAt, sellMfiAt, buyMfiAt float64) *PercentBandMFIStrategy {
	return &PercentBandMFIStrategy{
		MoneyFlowIndex: volume.NewMfi[float64](),
		PercentB:       volatility.NewPercentB[float64](),
		SellPercentBAt: sellPercentBAt,
		BuyPercentBAt:  buyPercentBAt,
		SellMfiAt:      sellMfiAt,
		BuyMfiAt:       buyMfiAt,
	}
}

// Name returns the name of the strategy.
func (m *PercentBandMFIStrategy) Name() string {
	return fmt.Sprintf("PercentB (%.2f,%.2f) and MFI Strategy (%.2f,%.2f)", m.SellPercentBAt, m.BuyPercentBAt, m.SellMfiAt, m.BuyMfiAt)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (m *PercentBandMFIStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 4)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[2]), 2)
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[3])

	mfis := m.MoneyFlowIndex.Compute(highs, lows, closings[0], volumes)
	mfis = helper.Shift(mfis, m.PercentB.IdlePeriod()-m.MoneyFlowIndex.IdlePeriod(), 0)
	pb := m.PercentB.Compute(closings[1])

	actions := helper.Operate(pb, mfis, func(b, mfi float64) strategy.Action {
		if b > m.BuyPercentBAt && mfi > m.BuyMfiAt {
			return strategy.Buy
		}

		if b < m.SellPercentBAt && mfi < m.SellMfiAt {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// strategy starts only after a full period.
	actions = helper.Shift(actions, m.PercentB.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (m *PercentBandMFIStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[0] -> closings
	//                 closings[1] -> money flow index
	//                 closings[2] -> percent b
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
		3,
	)
	volumes := asset.SnapshotsAsVolumes(snapshots[4])

	mfis := m.MoneyFlowIndex.Compute(highs, lows, closingsSplice[0], volumes)
	mfis = helper.Shift(mfis, m.PercentB.IdlePeriod()-m.MoneyFlowIndex.IdlePeriod(), 0)
	pb := m.PercentB.Compute(closingsSplice[2])

	closingsSplice[1] = helper.Skip(closingsSplice[1], m.PercentB.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(m, snapshots[5])
	actions = helper.Skip(actions, m.PercentB.IdlePeriod())
	outcomes = helper.Skip(outcomes, m.PercentB.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(m.Name(), dates)
	report.AddChart()
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Money Flow Index", mfis), 1)
	report.AddColumn(helper.NewNumericReportColumn("%B", pb), 2)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1, 2)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 3)

	return report
}
