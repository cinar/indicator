// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volatility"
)

// DonchianChannelBreakoutStrategy demonstrates how to compose Donchian Channels
// into an illustrative channel breakout strategy.
type DonchianChannelBreakoutStrategy struct {
	// DonchianChannel represents the configuration parameters for calculating the Donchian Channel.
	DonchianChannel *volatility.DonchianChannel[float64]
}

// NewDonchianChannelBreakoutStrategy initializes an example DonchianChannelBreakoutStrategy instance with default parameters.
func NewDonchianChannelBreakoutStrategy() *DonchianChannelBreakoutStrategy {
	return &DonchianChannelBreakoutStrategy{
		DonchianChannel: volatility.NewDonchianChannel[float64](),
	}
}

// Name returns the name of the example strategy.
func (*DonchianChannelBreakoutStrategy) Name() string {
	return "Donchian Channel Breakout Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (d *DonchianChannelBreakoutStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 3)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[2])

	uppers, middles, lowers := d.DonchianChannel.ComputeWithContext(ctx, highs, lows)
	go helper.DrainWithContext(ctx, middles)

	closings = helper.SkipWithContext(ctx, closings, d.DonchianChannel.IdlePeriod())

	actions := helper.Operate3WithContext(ctx, uppers, lowers, closings, func(upper, lower, closing float64) strategy.Action {
		if closing >= upper {
			return strategy.Buy
		}

		if closing <= lower {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Donchian Channel starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, d.DonchianChannel.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (d *DonchianChannelBreakoutStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs   -|
	// snapshots[2] -> lows    -+-> DonchianChannel.Compute -> upper, middle, lower
	// snapshots[3] -> closings -> close
	// snapshots[4] -> actions  -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := asset.SnapshotsAsClosings(snapshots[3])

	uppers, middles, lowers := d.DonchianChannel.Compute(highs, lows)
	uppers = helper.Shift(uppers, d.DonchianChannel.IdlePeriod(), 0)
	middles = helper.Shift(middles, d.DonchianChannel.IdlePeriod(), 0)
	lowers = helper.Shift(lowers, d.DonchianChannel.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(d, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(d.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("Upper", uppers))
	report.AddColumn(helper.NewNumericReportColumn("Middle", middles))
	report.AddColumn(helper.NewNumericReportColumn("Lower", lowers))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (d *DonchianChannelBreakoutStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return d.ComputeWithContext(context.Background(), snapshots)
}
