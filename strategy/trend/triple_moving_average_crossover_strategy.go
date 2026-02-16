// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultTripleMovingAverageCrossoverStrategyFastPeriod is the default triple moving average crossover strategy fast period.
	DefaultTripleMovingAverageCrossoverStrategyFastPeriod = 21

	// DefaultTripleMovingAverageCrossoverStrategyMediumPeriod is the default triple moving average crossover strategy medium period.
	DefaultTripleMovingAverageCrossoverStrategyMediumPeriod = 50

	// DefaultTripleMovingAverageCrossoverStrategySlowPeriod is the default triple moving average crossover strategy slow period.
	DefaultTripleMovingAverageCrossoverStrategySlowPeriod = 200
)

// TripleMovingAverageCrossoverStrategy defines the parameters used to calculate the Triple Moving Average Crossover
// trading strategy. This strategy uses three Exponential Moving Averages (EMAs) with different lengths to identify
// potential buy and sell signals.
// - A buy signal is generated when the **fastest** EMA crosses above both the **medium** and **slowest** EMAs.
// - A sell signal is generated when the fastest EMA crosses below both the medium and slowest EMAs.
// - Otherwise, the strategy recommends holding the asset.
type TripleMovingAverageCrossoverStrategy struct {
	// FastEma is the fastest EMA.
	FastEma *trend.Ema[float64]

	// MediumEma is the meium EMA.
	MediumEma *trend.Ema[float64]

	// SlowEma is the slowest EMA.
	SlowEma *trend.Ema[float64]
}

// NewTripleMovingAverageCrossoverStrategy function initializes a new Triple Moving Average Crossover strategy instance with the default parameters.
func NewTripleMovingAverageCrossoverStrategy() *TripleMovingAverageCrossoverStrategy {
	return NewTripleMovingAverageCrossoverStrategyWith(
		DefaultTripleMovingAverageCrossoverStrategyFastPeriod,
		DefaultTripleMovingAverageCrossoverStrategyMediumPeriod,
		DefaultTripleMovingAverageCrossoverStrategySlowPeriod,
	)
}

// NewTripleMovingAverageCrossoverStrategyWith function initializes a new Triple Moving Average Crossover strategy instance with the given periods.
func NewTripleMovingAverageCrossoverStrategyWith(fastPeriod, mediumPeriod, slowPeriod int) *TripleMovingAverageCrossoverStrategy {
	return &TripleMovingAverageCrossoverStrategy{
		FastEma:   trend.NewEmaWithPeriod[float64](fastPeriod),
		MediumEma: trend.NewEmaWithPeriod[float64](mediumPeriod),
		SlowEma:   trend.NewEmaWithPeriod[float64](slowPeriod),
	}
}

// Name returns the name of the strategy.
func (*TripleMovingAverageCrossoverStrategy) Name() string {
	return "Triple Moving Average Crossover Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (t *TripleMovingAverageCrossoverStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	fastEmas, mediumEmas, slowEmas := t.calculateEmas(c)

	actions := helper.Operate3(fastEmas, mediumEmas, slowEmas, func(fastEma, mediumEma, slowEma float64) strategy.Action {
		// A buy signal is generated when the **fastest** EMA crosses above both the **medium** and **slowest** EMAs.
		if (fastEma > mediumEma) && (fastEma > slowEma) {
			return strategy.Buy
		}

		// A sell signal is generated when the fastest EMA crosses below both the medium and slowest EMAs.
		if (fastEma < mediumEma) && (fastEma < slowEma) {
			return strategy.Sell
		}

		// Otherwise, the strategy recommends holding the asset.
		return strategy.Hold
	})

	// Generate a Hold signal during the idle period.
	actions = helper.Shift(actions, t.SlowEma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (t *TripleMovingAverageCrossoverStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings
	// snapshots[2] -> fastEmas
	//                 mediumEmas
	//                 slowEmas
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshots[0]),
		t.SlowEma.IdlePeriod(),
	)

	closingsSplice := helper.Duplicate(
		helper.Skip(
			asset.SnapshotsAsClosings(snapshots[1]),
			t.SlowEma.IdlePeriod(),
		),
		2,
	)

	fastEmas, mediumEmas, slowEmas := t.calculateEmas(snapshots[2])

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshots[3])

	annotations := helper.Skip(
		strategy.ActionsToAnnotations(actions),
		t.SlowEma.IdlePeriod(),
	)

	outcomes = helper.MultiplyBy(
		helper.Skip(
			outcomes,
			t.SlowEma.IdlePeriod(),
		),
		100,
	)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[0]))

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]), 1)
	report.AddColumn(helper.NewNumericReportColumn("Fast", fastEmas), 1)
	report.AddColumn(helper.NewNumericReportColumn("Medium", mediumEmas), 1)
	report.AddColumn(helper.NewNumericReportColumn("Slow", slowEmas), 1)

	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)
	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// calculateEmas calculates the fast, medium, and slow EMAs.
func (t *TripleMovingAverageCrossoverStrategy) calculateEmas(c <-chan *asset.Snapshot) (<-chan float64, <-chan float64, <-chan float64) {
	closings := helper.Duplicate(asset.SnapshotsAsClosings(c), 3)

	fastEmas := helper.Skip(
		t.FastEma.Compute(closings[0]),
		t.SlowEma.IdlePeriod()-t.FastEma.IdlePeriod(),
	)

	mediumEmas := helper.Skip(
		t.MediumEma.Compute(closings[1]),
		t.SlowEma.IdlePeriod()-t.MediumEma.IdlePeriod(),
	)

	slowEmas := t.SlowEma.Compute(closings[2])

	return fastEmas, mediumEmas, slowEmas
}
