// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

// MacdSignalMode selects how MacdStrategy turns a MACD/signal-line crossover into an action.
type MacdSignalMode int

const (
	// LevelTriggered evaluates the crossover condition independently on every bar, so it keeps
	// returning Buy (or Sell) on every consecutive bar for which the condition holds, not just
	// the bar on which the crossing actually happened. This is the existing default behavior.
	LevelTriggered MacdSignalMode = iota

	// EdgeTriggered fires Buy/Sell only once, on the bar where the MACD/signal-line crossing
	// actually occurs, by comparing the previous bar's MACD/signal pair against the current
	// one. This requires one additional bar of history, so the strategy's idle period is one
	// bar longer than in LevelTriggered mode.
	EdgeTriggered
)

// MacdStrategy demonstrates how to compose the Moving Average Convergence Divergence (MACD)
// and its signal line into an illustrative, zero-line-filtered MACD crossover strategy. Buy
// signals fire only on a MACD-above-signal crossing that occurs while MACD is still below
// zero; Sell signals fire only on a MACD-below-signal crossing while MACD is still above
// zero.
type MacdStrategy struct {
	// Macd represents the configuration parameters for calculating the
	// Moving Average Convergence Divergence (MACD).
	Macd *trend.Macd[float64]

	// SignalMode selects between LevelTriggered (default) and EdgeTriggered crossover
	// detection. See MacdSignalMode for details.
	SignalMode MacdSignalMode
}

// NewMacdStrategy initializes an example MacdStrategy instance with default parameters.
func NewMacdStrategy() *MacdStrategy {
	return NewMacdStrategyWith(
		trend.DefaultMacdPeriod1,
		trend.DefaultMacdPeriod2,
		trend.DefaultMacdPeriod3,
	)
}

// NewMacdStrategyWith initializes an example MacdStrategyWith instance with default parameters.
func NewMacdStrategyWith(period1, period2, period3 int) *MacdStrategy {
	return &MacdStrategy{
		Macd: trend.NewMacdWithPeriod[float64](
			period1,
			period2,
			period3,
		),
	}
}

// Name returns the name of the example strategy.
func (m *MacdStrategy) Name() string {
	return fmt.Sprintf("MACD Strategy (%d,%d,%d)",
		m.Macd.Ema1.Period,
		m.Macd.Ema2.Period,
		m.Macd.Ema3.Period,
	)
}

// ComputeWithContext processes the provided asset snapshots and generates an
// illustrative stream of actions.
func (m *MacdStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots)

	macds, signals := m.Macd.ComputeWithContext(ctx, closings)

	if m.SignalMode == EdgeTriggered {
		return m.computeEdgeTriggered(ctx, macds, signals)
	}

	actions := helper.OperateWithContext(ctx, macds, signals, func(macd, signal float64) strategy.Action {
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
	actions = helper.ShiftWithContext(ctx, actions, m.Macd.IdlePeriod(), strategy.Hold)

	return actions
}

// computeEdgeTriggered implements the EdgeTriggered SignalMode: it fires Buy/Sell only once,
// on the bar where the MACD/signal-line crossing occurs, by pairing each bar's MACD/signal
// values with the previous bar's.
func (m *MacdStrategy) computeEdgeTriggered(ctx context.Context, macds, signals <-chan float64) <-chan strategy.Action {
	macds = helper.BufferedWithContext(ctx, macds, 2)
	signals = helper.BufferedWithContext(ctx, signals, 2)

	macdInputs := helper.DuplicateWithContext(ctx, macds, 2)
	macdInputs[1] = helper.SkipWithContext(ctx, macdInputs[1], 1)

	signalInputs := helper.DuplicateWithContext(ctx, signals, 2)
	signalInputs[1] = helper.SkipWithContext(ctx, signalInputs[1], 1)

	actions := helper.Operate4WithContext(ctx, macdInputs[0], signalInputs[0], macdInputs[1], signalInputs[1], func(prevMacd, prevSignal, macd, signal float64) strategy.Action {
		// A MACD value crossing above the signal line while still below zero suggests a bullish trend.
		if (macd > signal) && (prevMacd <= prevSignal) && (macd < 0) {
			return strategy.Buy
		}

		// A MACD value crossing below the signal line while still above zero suggests a bearish trend.
		if (macd < signal) && (prevMacd >= prevSignal) && (macd > 0) {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// MACD starts only after a full period, plus one more bar consumed by the
	// previous/current pairing used for crossover edge-detection.
	actions = helper.ShiftWithContext(ctx, actions, m.Macd.IdlePeriod()+1, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (m *MacdStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return m.ComputeWithContext(context.Background(), snapshots)
}
