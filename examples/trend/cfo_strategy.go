// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

// CfoStrategy demonstrates how to compose the Chande Forecast Oscillator (CFO)
// indicator into an illustrative zero-line crossover strategy.
type CfoStrategy struct {
	// Cfo represents the configuration parameters for calculating the
	// Chande Forecast Oscillator (CFO).
	Cfo *trend.Cfo[float64]
}

// NewCfoStrategy initializes an example CfoStrategy instance with default parameters.
func NewCfoStrategy() *CfoStrategy {
	return &CfoStrategy{
		Cfo: trend.NewCfo[float64](),
	}
}

// Name returns the name of the example strategy.
func (*CfoStrategy) Name() string {
	return "Cfo Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an
// illustrative stream of actions.
func (c *CfoStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots)

	cfo := c.Cfo.ComputeWithContext(ctx, closings)
	cfo = helper.BufferedWithContext(ctx, cfo, 2)

	inputs := helper.DuplicateWithContext(ctx, cfo, 2)

	// Skip the first value
	inputs[1] = helper.SkipWithContext(ctx, inputs[1], 1)

	actions := helper.OperateWithContext(ctx, inputs[0], inputs[1], func(b, c float64) strategy.Action {
		// A CFO value crossing above zero suggests a bullish trend.
		if c >= 0 && b < 0 {
			return strategy.Buy
		}

		// A CFO value crossing below zero indicates a bearish trend.
		if c <= 0 && b > 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// CFO starts only after the period.
	actions = helper.ShiftWithContext(ctx, actions, c.Cfo.Mlr.IdlePeriod()+1, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
func (c *CfoStrategy) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	dates := asset.SnapshotsAsDates(snapshotsSplice[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[2]), 2)
	cfo := helper.Shift(c.Cfo.Compute(closings[1]), c.Cfo.Mlr.IdlePeriod()+1, 0)

	actions, outcomes := strategy.ComputeWithOutcome(c, snapshotsSplice[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(c.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("CFO", cfo), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (c *CfoStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return c.ComputeWithContext(context.Background(), snapshots)
}
