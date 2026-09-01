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

// KamaStrategy demonstrates how to compose Kaufman's Adaptive Moving Averages (KAMAs)
// with fast and slow periods into an illustrative adaptive moving average crossover strategy.
type KamaStrategy struct {
	// Kama represents the configuration parameters for calculating the Kaufman's Adaptive Moving Average (KAMA).
	Kama *trend.Kama[float64]
}

// NewKamaStrategy initializes an example KamaStrategy instance with default parameters.
func NewKamaStrategy() *KamaStrategy {
	return NewKamaStrategyWith(
		trend.DefaultKamaErPeriod,
		trend.DefaultKamaFastScPeriod,
		trend.DefaultKamaSlowScPeriod,
	)
}

// NewKamaStrategyWith initializes an example KamaStrategyWith instance with default parameters.
func NewKamaStrategyWith(erPeriod, fastScPeriod, slowScPeriod int) *KamaStrategy {
	return &KamaStrategy{
		Kama: trend.NewKamaWith[float64](
			erPeriod,
			fastScPeriod,
			slowScPeriod,
		),
	}
}

// Name returns the name of the example strategy.
func (k *KamaStrategy) Name() string {
	return k.Kama.String()
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (k *KamaStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.DuplicateWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshots), 2)
	closingsSplice[1] = helper.SkipWithContext(ctx, closingsSplice[1], k.Kama.IdlePeriod())

	kamas := k.Kama.ComputeWithContext(ctx, closingsSplice[0])

	actions := helper.OperateWithContext(ctx, kamas, closingsSplice[1], func(kama, closing float64) strategy.Action {
		// A closing price crossing above the KAMA suggests a bullish trend.
		if closing > kama {
			return strategy.Buy
		}

		// While crossing below the KAMA indicates a bearish trend.
		if closing < kama {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// KAMA starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, k.Kama.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an illustrative report annotated with example actions.
func (k *KamaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> kama
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(c, 3)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		k.Kama.IdlePeriod(),
	)

	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[1]), 2)
	closingsSplice[1] = helper.Skip(closingsSplice[1], k.Kama.IdlePeriod())

	kamas := k.Kama.Compute(closingsSplice[0])

	actions, outcomes := strategy.ComputeWithOutcome(k, snapshotsSplice[2])
	actions = helper.Skip(actions, k.Kama.IdlePeriod())
	outcomes = helper.Skip(outcomes, k.Kama.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(k.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("KAMA", kamas), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (k *KamaStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return k.ComputeWithContext(context.Background(), snapshots)
}
