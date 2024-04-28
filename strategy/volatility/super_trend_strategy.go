// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volatility"
)

// SuperTrendStrategy represents the configuration parameters for calculating the Super Trend strategy. A closing
// value crossing above the Super Trend suggets a Buy signal, while crossing below the Super Trend indivates a
// Sell signal.
type SuperTrendStrategy struct {
	// SuperTrend represents the configuration parameters for calculating the Super Trend.
	SuperTrend *volatility.SuperTrend[float64]
}

// NewSuperTrendStrategy function initializes a new Super Trend strategy instance.
func NewSuperTrendStrategy() *SuperTrendStrategy {
	return NewSuperTrendStrategyWith(volatility.NewSuperTrend[float64]())
}

// NewSuperTrendStrategyWith function initializes a new Super Trend strategy with the given Super Trend instance.
func NewSuperTrendStrategyWith(superTrend *volatility.SuperTrend[float64]) *SuperTrendStrategy {
	return &SuperTrendStrategy{
		SuperTrend: superTrend,
	}
}

// Name returns the name of the strategy.
func (s *SuperTrendStrategy) Name() string {
	return fmt.Sprintf("Super Trend Strategy (%s, %.1f)", s.SuperTrend.Atr.Ma.String(), s.SuperTrend.Multiplier)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *SuperTrendStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshotsSplice[2]),
		2,
	)

	superTrends := s.SuperTrend.Compute(highs, lows, closingsSplice[0])

	closingsSplice[1] = helper.Skip(closingsSplice[1], s.SuperTrend.IdlePeriod())

	actions := helper.Operate(superTrends, closingsSplice[1], func(superTrend, closing float64) strategy.Action {
		if superTrend < closing {
			return strategy.Buy
		}

		if superTrend > closing {
			return strategy.Sell
		}

		return strategy.Hold
	})

	actions = strategy.NormalizeActions(actions)

	// Super Trend starts only after a full period.
	actions = helper.Shift(actions, s.SuperTrend.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
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
