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

const (
	// DefaultTrimaStrategyShortPeriod is the first TRIMA period.
	DefaultTrimaStrategyShortPeriod = 20

	// DefaultTrimaStrategyLongPeriod is the second TRIMA period.
	DefaultTrimaStrategyLongPeriod = 50
)

// TrimaStrategy demonstrates how to compose Triangular Moving Averages (TRIMAs)
// with fast and slow periods into an illustrative moving average crossover strategy.
type TrimaStrategy struct {
	// Trima1 represents the configuration parameters for calculating the short TRIMA.
	Short *trend.Trima[float64]

	// Trima2 represents the configuration parameters for calculating the long TRIMA.
	Long *trend.Trima[float64]
}

// NewTrimaStrategy initializes an example TrimaStrategy instance with default parameters.
// with the default parameters.
func NewTrimaStrategy() *TrimaStrategy {
	t := &TrimaStrategy{
		Short: trend.NewTrima[float64](),
		Long:  trend.NewTrima[float64](),
	}

	t.Short.Period = DefaultTrimaStrategyShortPeriod
	t.Long.Period = DefaultTrimaStrategyLongPeriod

	return t
}

// Name returns the name of the example strategy.
func (*TrimaStrategy) Name() string {
	return "TRIMA Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an
// illustrative stream of actions.
func (t *TrimaStrategy) ComputeWithContext(ctx context.Context, c <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := helper.DuplicateWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, c), 2)

	shorts := t.Short.ComputeWithContext(ctx, closings[0])
	longs := t.Long.ComputeWithContext(ctx, closings[1])

	shorts = helper.SkipWithContext(ctx, shorts, t.Long.IdlePeriod()-t.Short.IdlePeriod())

	actions := helper.OperateWithContext(ctx, shorts, longs, func(short, long float64) strategy.Action {
		if short > long {
			return strategy.Buy
		}

		if long > short {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// TRIMA starts only after a full periods for each EMA used.
	actions = helper.ShiftWithContext(ctx, actions, t.Long.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
func (t *TrimaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> shorts
	//                 closings[1] -> longs
	//                 closings[2] -> closings
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 3)

	shorts := t.Short.Compute(closings[0])
	longs := t.Long.Compute(closings[1])

	shorts = helper.Skip(shorts, t.Long.IdlePeriod()-t.Short.IdlePeriod())
	shorts = helper.Shift(shorts, t.Long.IdlePeriod(), 0)
	longs = helper.Shift(longs, t.Long.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[2]))
	report.AddColumn(helper.NewNumericReportColumn("Short", shorts), 1)
	report.AddColumn(helper.NewNumericReportColumn("Long", longs), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (t *TrimaStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	return t.ComputeWithContext(context.Background(), c)
}
