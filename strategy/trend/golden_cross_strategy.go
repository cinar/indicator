// Copyright (c) 2021-2024 Onur Cinar.
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
	// DefaultGoldenCrossStrategyFastPeriod is the default golden cross strategy fast period.
	DefaultGoldenCrossStrategyFastPeriod = 21

	// DefaultGoldenCrossStrategyMediumPeriod is the default golden cross strategy medium period.
	DefaultGoldenCrossStrategyMediumPeriod = 50

	// DefaultGoldenCrossStrategySlowPeriod is the default golden cross strategy slow period.
	DefaultGoldenCrossStrategySlowPeriod = 200
)

// GoldenCrossStrategy defines the parameters used to calculate the Golden Cross trading strategy. This strategy uses
// three Exponential Moving Averages (EMAs) with different lengths to identify potential buy and sell signals.
// - A buy signal is generated when the **fastest** EMA crosses above both the **medium** and **slowest** EMAs.
// - A sell signal is generated when the fastest EMA crosses below both the medium and slowest EMAs.
// - Otherwise, the strategy recommends holding the asset.
type GoldenCrossStrategy struct {
	// FastEma is the fastest EMA.
	FastEma *trend.Ema[float64]

	// MediumEma is the meium EMA.
	MediumEma *trend.Ema[float64]

	// SlowEma is the slowest EMA.
	SlowEma *trend.Ema[float64]
}

// NewGoldenCrossStrategy function initializes a new Golden Cross strategy instance with the default parameters.
func NewGoldenCrossStrategy() *GoldenCrossStrategy {
	return NewGoldenCrossStrategyWith(
		DefaultGoldenCrossStrategyFastPeriod,
		DefaultGoldenCrossStrategyMediumPeriod,
		DefaultGoldenCrossStrategySlowPeriod,
	)
}

// NewGoldenCrossStrategyWith function initializes a new Golden Cross strategy instance with the given periods.
func NewGoldenCrossStrategyWith(fastPeriod, mediumPeriod, slowPeriod int) *GoldenCrossStrategy {
	return &GoldenCrossStrategy{
		FastEma:   trend.NewEmaWithPeriod[float64](fastPeriod),
		MediumEma: trend.NewEmaWithPeriod[float64](mediumPeriod),
		SlowEma:   trend.NewEmaWithPeriod[float64](slowPeriod),
	}
}

// Name returns the name of the strategy.
func (*GoldenCrossStrategy) Name() string {
	return "Golden Cross Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (g *GoldenCrossStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	fastEmas, mediumEmas, slowEmas := g.calculateEmas(c)

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

	// Normalize actions
	actions = strategy.NormalizeActions(actions)

	// Generate a Hold signal during the idle period.
	actions = helper.Shift(actions, g.SlowEma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (g *GoldenCrossStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
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
		g.SlowEma.IdlePeriod(),
	)

	closingsSplice := helper.Duplicate(
		helper.Skip(
			asset.SnapshotsAsClosings(snapshots[1]),
			g.SlowEma.IdlePeriod(),
		),
		2,
	)

	fastEmas, mediumEmas, slowEmas := g.calculateEmas(snapshots[2])

	actions, outcomes := strategy.ComputeWithOutcome(g, snapshots[3])

	annotations := helper.Skip(
		strategy.ActionsToAnnotations(actions),
		g.SlowEma.IdlePeriod(),
	)

	outcomes = helper.MultiplyBy(
		helper.Skip(
			outcomes,
			g.SlowEma.IdlePeriod(),
		),
		100,
	)

	report := helper.NewReport(g.Name(), dates)
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
func (g *GoldenCrossStrategy) calculateEmas(c <-chan *asset.Snapshot) (<-chan float64, <-chan float64, <-chan float64) {
	closings := helper.Duplicate(asset.SnapshotsAsClosings(c), 3)

	fastEmas := helper.Skip(
		g.FastEma.Compute(closings[0]),
		g.SlowEma.IdlePeriod()-g.FastEma.IdlePeriod(),
	)

	mediumEmas := helper.Skip(
		g.MediumEma.Compute(closings[1]),
		g.SlowEma.IdlePeriod()-g.MediumEma.IdlePeriod(),
	)

	slowEmas := g.SlowEma.Compute(closings[2])

	return fastEmas, mediumEmas, slowEmas
}
