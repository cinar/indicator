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

const (
	// DefaultTsiStrategySignalPeriod is the default signal line period of 12.
	DefaultTsiStrategySignalPeriod = 12
)

// TsiStrategy demonstrates how to compose the True Strength Index (TSI) indicator
// and its signal line into an illustrative crossover strategy.
type TsiStrategy struct {
	// Tsi represents the configuration parameters for calculating the True Strength Index (TSI).
	Tsi *trend.Tsi[float64]

	// Signal line is the moving average of the TSI.
	Signal trend.Ma[float64]
}

// NewTsiStrategy initializes an example TsiStrategy instance with default parameters.
func NewTsiStrategy() *TsiStrategy {
	return NewTsiStrategyWith(
		trend.DefaultTsiFirstSmoothingPeriod,
		trend.DefaultTsiSecondSmoothingPeriod,
		DefaultTsiStrategySignalPeriod,
	)
}

// NewTsiStrategyWith initializes an example TsiStrategyWith instance with default parameters.
func NewTsiStrategyWith(firstSmoothingPeriod, secondSmoothingPeriod, signalPeriod int) *TsiStrategy {
	return &TsiStrategy{
		Tsi: trend.NewTsiWith[float64](
			firstSmoothingPeriod,
			secondSmoothingPeriod,
		),

		Signal: trend.NewEmaWithPeriod[float64](signalPeriod),
	}
}

// Name returns the name of the example strategy.
func (t *TsiStrategy) Name() string {
	return fmt.Sprintf("Tsi Strategy (%s,%s)",
		t.Tsi.String(),
		t.Signal.String(),
	)
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (t *TsiStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots)

	tsisSplice := helper.DuplicateWithContext(ctx, t.Tsi.ComputeWithContext(ctx, closings), 2)

	tsisSplice[0] = helper.SkipWithContext(ctx, tsisSplice[0], t.Signal.IdlePeriod())
	signals := trend.ComputeMaWithContext(ctx, t.Signal, tsisSplice[1])

	actions := helper.OperateWithContext(ctx, tsisSplice[0], signals, func(tsi, signal float64) strategy.Action {
		// When the TSI is above zero and crossing above the signal line suggests a bullish trend.
		if (tsi > 0) && (tsi > signal) {
			return strategy.Buy
		}

		// While TSI being below zero and crossing below the signal line indicates a bearish trend.
		if (tsi < 0) && (tsi < signal) {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// TSI and signal line start only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, t.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (t *TsiStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> tsi[0] -> tsi
	//                             -> tsi[1] -> signal
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(c, 3)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		t.IdlePeriod(),
	)

	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[1]), 2)
	closingsSplice[1] = helper.Skip(closingsSplice[1], t.IdlePeriod())

	tsisSplice := helper.Duplicate(t.Tsi.Compute(closingsSplice[0]), 2)
	tsisSplice[0] = helper.Skip(tsisSplice[0], t.Signal.IdlePeriod())

	signals := t.Signal.Compute(tsisSplice[1])

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshotsSplice[2])
	actions = helper.Skip(actions, t.IdlePeriod())
	outcomes = helper.Skip(outcomes, t.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))

	report.AddColumn(helper.NewNumericReportColumn("TSI", tsisSplice[0]), 1)
	report.AddColumn(helper.NewNumericReportColumn("Signal", signals), 1)

	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// IdlePeriod is the initial period that TSI strategy yield any results.
func (t *TsiStrategy) IdlePeriod() int {
	return t.Tsi.IdlePeriod() + t.Signal.IdlePeriod()
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (t *TsiStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return t.ComputeWithContext(context.Background(), snapshots)
}
