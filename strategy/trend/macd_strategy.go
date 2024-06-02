// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

// MacdStrategy represents the configuration parameters for calculating the
// MACD strategy. A MACD value crossing above the signal line suggests a
// bullish trend, while crossing below the signal line indicates a
// bearish trend.
type MacdStrategy struct {
	strategy.Strategy

	// Macd represents the configuration parameters for calculating the
	// Moving Average Convergence Divergence (MACD).
	Macd *trend.Macd[float64]
}

// NewMacdStrategy function initializes a new MACD strategy instance.
func NewMacdStrategy() *MacdStrategy {
	return NewMacdStrategyWith(
		trend.DefaultMacdPeriod1,
		trend.DefaultMacdPeriod2,
		trend.DefaultMacdPeriod3,
	)
}

// NewMacdStrategyWith function initializes a new MACD strategy instance with the given parameters.
func NewMacdStrategyWith(period1, period2, period3 int) *MacdStrategy {
	return &MacdStrategy{
		Macd: trend.NewMacdWithPeriod[float64](
			period1,
			period2,
			period3,
		),
	}
}

// Name returns the name of the strategy.
func (m *MacdStrategy) Name() string {
	return fmt.Sprintf("MACD Strategy (%d,%d,%d)",
		m.Macd.Ema1.Period,
		m.Macd.Ema2.Period,
		m.Macd.Ema3.Period,
	)
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (m *MacdStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosings(snapshots)

	macds, signals := m.Macd.Compute(closings)

	actions := helper.Operate(macds, signals, func(macd, signal float64) strategy.Action {
		// A MACD value crossing above signal line suggests a bullish trend.
		if (macd > signal) && (macd < 0) {
			return strategy.Buy
		}

		// A MACD value crossing below signal line suggests a bearish trend.
		if (signal > macd) && (macd > 0) {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// MACD starts only after a full period.
	actions = helper.Shift(actions, m.Macd.IdlePeriod(), strategy.Hold)

	actions = strategy.NormalizeActions(actions)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (m *MacdStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> macds, signals
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 2)

	macds, signals := m.Macd.Compute(closings[0])
	macds = helper.Shift(macds, m.Macd.IdlePeriod(), 0)
	signals = helper.Shift(signals, m.Macd.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(m, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(m.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("MACD", macds), 1)
	report.AddColumn(helper.NewNumericReportColumn("Signal", signals), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
