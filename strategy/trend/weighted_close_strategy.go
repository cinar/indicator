// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultWeightedCloseStrategyMaPeriod is the default Moving Average period of 20.
	DefaultWeightedCloseStrategyMaPeriod = 20
)

// WeightedCloseStrategy represents the configuration parameters for calculating
// the Weighted Close strategy. A weighted close crossing above the moving
// average suggests a bullish trend, while crossing below the moving
// average indicates a bearish trend.
type WeightedCloseStrategy struct {
	// WeightedClose represents the configuration parameters for calculating the weighted close.
	WeightedClose *trend.WeightedClose[float64]

	// Ma represents the configuration parameters for calculating the moving average.
	Ma trend.Ma[float64]
}

// NewWeightedCloseStrategy function initializes a new Weighted Close strategy instance.
func NewWeightedCloseStrategy() *WeightedCloseStrategy {
	return NewWeightedCloseStrategyWith(
		DefaultWeightedCloseStrategyMaPeriod,
	)
}

// NewWeightedCloseStrategyWith function initializes a new Weighted Close strategy instance
// with the given parameters.
func NewWeightedCloseStrategyWith(maPeriod int) *WeightedCloseStrategy {
	return &WeightedCloseStrategy{
		WeightedClose: trend.NewWeightedClose[float64](),
		Ma:            trend.NewSmaWithPeriod[float64](maPeriod),
	}
}

// Name returns the name of the strategy.
func (w *WeightedCloseStrategy) Name() string {
	return fmt.Sprintf("Weighted Close Strategy (%d)",
		w.Ma.IdlePeriod()+1,
	)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (w *WeightedCloseStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])

	wcSplice := helper.Duplicate(
		w.WeightedClose.Compute(highs, lows, closings),
		2,
	)

	mas := w.Ma.Compute(wcSplice[1])

	wcSplice[0] = helper.Skip(wcSplice[0], w.Ma.IdlePeriod())

	actions := helper.Operate(wcSplice[0], mas, func(wc, ma float64) strategy.Action {
		// A weighted close crossing above the moving average suggests a bullish trend.
		if wc > ma {
			return strategy.Buy
		}

		// A crossing below the moving average indicates a bearish trend.
		return strategy.Sell
	})

	// SMMA strategy starts only after a full period.
	actions = helper.Shift(actions, w.Ma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (w *WeightedCloseStrategy) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs
	// snapshots[2] -> lows
	// snapshots[3] -> closings[0] -> closings
	//                 closings[1] -> weighted closes[0] -> weighted closes
	//                                weighted closes[1] -> moving average
	// snapshots[4] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(snapshots, 5)

	dates := asset.SnapshotsAsDates(snapshotsSplice[0])
	highs := asset.SnapshotsAsHighs(snapshotsSplice[1])
	lows := asset.SnapshotsAsLows(snapshotsSplice[2])
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshotsSplice[3]),
		2,
	)

	wcSplice := helper.Duplicate(
		w.WeightedClose.Compute(highs, lows, closingsSplice[1]),
		2,
	)

	mas := w.Ma.Compute(wcSplice[1])

	actions, outcomes := strategy.ComputeWithOutcome(w, snapshotsSplice[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	dates = helper.Skip(dates, w.Ma.IdlePeriod())
	closingsSplice[0] = helper.Skip(closingsSplice[0], w.Ma.IdlePeriod())
	wcSplice[0] = helper.Skip(wcSplice[0], w.Ma.IdlePeriod())
	annotations = helper.Skip(annotations, w.Ma.IdlePeriod())
	outcomes = helper.Skip(outcomes, w.Ma.IdlePeriod())

	report := helper.NewReport(w.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[0]))
	report.AddColumn(helper.NewNumericReportColumn("Weighted Close", wcSplice[0]))
	report.AddColumn(helper.NewNumericReportColumn("Moving Average", mas))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
