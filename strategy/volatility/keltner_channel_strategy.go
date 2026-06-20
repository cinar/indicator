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

// KeltnerChannelStrategy represents the configuration parameters for calculating the Keltner Channel strategy.
// A closing above the upper band suggests a Sell signal, while a closing below the lower band suggests a Buy signal.
type KeltnerChannelStrategy struct {
	// KeltnerChannel represents the configuration parameters for calculating the Keltner Channel.
	KeltnerChannel *volatility.KeltnerChannel[float64]
}

// NewKeltnerChannelStrategy function initializes a new Keltner Channel strategy instance.
func NewKeltnerChannelStrategy() *KeltnerChannelStrategy {
	return &KeltnerChannelStrategy{
		KeltnerChannel: volatility.NewKeltnerChannel[float64](),
	}
}

// Name returns the name of the strategy.
func (*KeltnerChannelStrategy) Name() string {
	return "Keltner Channel Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates a stream of actionable recommendations.
func (k *KeltnerChannelStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 4)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[2])

	uppers, middles, lowers := k.KeltnerChannel.ComputeWithContext(ctx, highs, lows, closings)
	go helper.DrainWithContext(ctx, middles)

	closings2 := helper.SkipWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[3]), k.KeltnerChannel.IdlePeriod())

	actions := helper.Operate3WithContext(ctx, uppers, lowers, closings2, func(upper, lower, closing float64) strategy.Action {
		if closing > upper {
			return strategy.Sell
		}

		if closing < lower {
			return strategy.Buy
		}

		return strategy.Hold
	})

	// Keltner Channel starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, k.KeltnerChannel.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (k *KeltnerChannelStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs   -|
	// snapshots[2] -> lows    -+-> KeltnerChannel.Compute -> upper, middle, lower
	// snapshots[3] -> closings-|
	//                 closings -> close
	// snapshots[4] -> actions  -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[3]), 2)

	uppers, middles, lowers := k.KeltnerChannel.Compute(highs, lows, closings[0])
	uppers = helper.Shift(uppers, k.KeltnerChannel.IdlePeriod(), 0)
	middles = helper.Shift(middles, k.KeltnerChannel.IdlePeriod(), 0)
	lowers = helper.Shift(lowers, k.KeltnerChannel.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(k, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(k.Name(), dates)
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
func (k *KeltnerChannelStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return k.ComputeWithContext(context.Background(), snapshots)
}
