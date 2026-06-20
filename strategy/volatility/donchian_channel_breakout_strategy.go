// Copyright (c) 2021-2026 Onur Cinar.
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

// DonchianChannelBreakoutStrategy represents the configuration parameters for calculating the Donchian Channel
// Breakout strategy. A closing at or above the upper channel suggests a Buy signal, while a closing at or below
// the lower channel suggests a Sell signal.
type DonchianChannelBreakoutStrategy struct {
	// DonchianChannel represents the configuration parameters for calculating the Donchian Channel.
	DonchianChannel *volatility.DonchianChannel[float64]
}

// NewDonchianChannelBreakoutStrategy function initializes a new Donchian Channel Breakout strategy instance.
func NewDonchianChannelBreakoutStrategy() *DonchianChannelBreakoutStrategy {
	return &DonchianChannelBreakoutStrategy{
		DonchianChannel: volatility.NewDonchianChannel[float64](),
	}
}

// Name returns the name of the strategy.
func (*DonchianChannelBreakoutStrategy) Name() string {
	return "Donchian Channel Breakout Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates a stream of actionable recommendations.
func (d *DonchianChannelBreakoutStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := helper.DuplicateWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshots),
		2,
	)

	uppers, middles, lowers := d.DonchianChannel.ComputeWithContext(ctx, closings[0])
	go helper.DrainWithContext(ctx, middles)

	closings[1] = helper.SkipWithContext(ctx, closings[1], d.DonchianChannel.IdlePeriod())

	actions := helper.Operate3WithContext(ctx, uppers, lowers, closings[1], func(upper, lower, closing float64) strategy.Action {
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

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (d *DonchianChannelBreakoutStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> upper
	//                             -> middle
	//                             -> lower
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 2)

	uppers, middles, lowers := d.DonchianChannel.Compute(closings[0])
	uppers = helper.Shift(uppers, d.DonchianChannel.IdlePeriod(), 0)
	middles = helper.Shift(middles, d.DonchianChannel.IdlePeriod(), 0)
	lowers = helper.Shift(lowers, d.DonchianChannel.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(d, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(d.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
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
