// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// IchimokuCloudStrategy demonstrates how to compose Tenkan-sen and Kijun-sen lines
// from the Ichimoku Cloud into an illustrative crossover strategy.
type IchimokuCloudStrategy struct {
	// IchimokuCloud represents the configuration parameters for calculating the Ichimoku Cloud.
	IchimokuCloud *momentum.IchimokuCloud[float64]
}

// NewIchimokuCloudStrategy initializes an example IchimokuCloudStrategy instance with default parameters.
func NewIchimokuCloudStrategy() *IchimokuCloudStrategy {
	return &IchimokuCloudStrategy{
		IchimokuCloud: momentum.NewIchimokuCloud[float64](),
	}
}

// Name returns the name of the example strategy.
func (*IchimokuCloudStrategy) Name() string {
	return "Ichimoku Cloud Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (i *IchimokuCloudStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 3)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[2])

	closingsSplice := helper.DuplicateWithContext(ctx, closings, 2)

	cl, bl, lsa, lsb, ll := i.IchimokuCloud.ComputeWithContext(ctx, highs, lows, closingsSplice[0])

	// Lagging line is not used in the core logic, drain it to prevent blocking
	go helper.DrainWithContext(ctx, ll)

	// cl, bl, lsa, and lsb are each IchimokuCloud.LaggingPeriod shorter than
	// IdlePeriod() alone would suggest, since the Chikou Span's forward lookahead
	// also truncates their tail. Skip the extra LaggingPeriod off closings so it
	// lines up 1:1 with them.
	alignPeriod := i.IchimokuCloud.IdlePeriod() + i.IchimokuCloud.LaggingPeriod

	actions := helper.Operate5WithContext(ctx, helper.SkipWithContext(ctx, closingsSplice[1], alignPeriod),
		cl,
		bl,
		lsa,
		lsb,
		func(c, conversion, base, spanA, spanB float64) strategy.Action {
			if c > spanA && c > spanB && conversion > base && spanA > spanB {
				return strategy.Buy
			}

			if c < spanA && c < spanB && conversion < base && spanA < spanB {
				return strategy.Sell
			}

			return strategy.Hold
		},
	)

	// Shift the actions to account for the idle period
	return helper.ShiftWithContext(ctx, actions, alignPeriod, strategy.Hold)
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (i *IchimokuCloudStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[2])
	highs := asset.SnapshotsAsHighs(snapshots[3])
	lows := asset.SnapshotsAsLows(snapshots[4])
	closingsForCloud := asset.SnapshotsAsClosings(snapshots[5])

	cl, bl, lsa, lsb, ll := i.IchimokuCloud.Compute(highs, lows, closingsForCloud)

	// Lagging line is not used in the report right now, drain it.
	go helper.Drain(ll)

	// See the matching comment in ComputeWithContext for why the extra LaggingPeriod is needed here.
	alignPeriod := i.IchimokuCloud.IdlePeriod() + i.IchimokuCloud.LaggingPeriod

	clShifted := helper.Shift(cl, alignPeriod, 0)
	blShifted := helper.Shift(bl, alignPeriod, 0)
	lsaShifted := helper.Shift(lsa, alignPeriod, 0)
	lsbShifted := helper.Shift(lsb, alignPeriod, 0)

	actions, outcomes := strategy.ComputeWithOutcome(i, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(i.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("Conversion Line", clShifted))
	report.AddColumn(helper.NewNumericReportColumn("Base Line", blShifted))
	report.AddColumn(helper.NewNumericReportColumn("Leading Span A", lsaShifted))
	report.AddColumn(helper.NewNumericReportColumn("Leading Span B", lsbShifted))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *IchimokuCloudStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return i.ComputeWithContext(context.Background(), snapshots)
}
