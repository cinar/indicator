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

// KdjStrategy represents the configuration parameters for calculating the KDJ strategy.
// Generates BUY action when j value crosses above both k and d values.
// Generates SELL action when j value crosses below both k and d values.
type KdjStrategy struct {
	strategy.Strategy

	// Kdj represents the configuration parameters for calculating the KDJ.
	Kdj *trend.Kdj[float64]
}

// NewKdjStrategy function initializes a new KDJ strategy instance.
func NewKdjStrategy() *KdjStrategy {
	return &KdjStrategy{
		Kdj: trend.NewKdj[float64](),
	}
}

// Name returns the name of the strategy.
func (*KdjStrategy) Name() string {
	return "KDJ Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (kdj *KdjStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshots := helper.Duplicate(c, 3)
	highs := asset.SnapshotsAsHighs(snapshots[0])
	lows := asset.SnapshotsAsLows(snapshots[1])
	closings := asset.SnapshotsAsClosings(snapshots[2])

	k, d, j := kdj.Kdj.Compute(highs, lows, closings)
	js := helper.Duplicate(j, 2)

	jk := helper.Subtract(js[0], k)
	jd := helper.Subtract(js[1], d)

	actions := helper.Operate(jk, jd, func(a, b float64) strategy.Action {
		// Generates BUY action when j value crosses above both k and d values.
		if a > 0 && b > 0 {
			return strategy.Buy
		}

		// Generates SELL action when j value crosses below both k and d values.
		if a < 0 && b < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	actions = strategy.NormalizeActions(actions)

	// KDJ starts only after a full period.
	actions = helper.Shift(actions, kdj.Kdj.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (kdj *KdjStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[1] |> kdj
	//                 closings[0] -> closings
	// snapshots[4] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsHighs(snapshots[2])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[3]), 2)

	k, d, j := kdj.Kdj.Compute(highs, lows, closings[1])
	k = helper.Shift(k, kdj.Kdj.IdlePeriod(), 0)
	d = helper.Shift(d, kdj.Kdj.IdlePeriod(), 0)
	j = helper.Shift(j, kdj.Kdj.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(kdj, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(kdj.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("K", k), 1)
	report.AddColumn(helper.NewNumericReportColumn("D", d), 1)
	report.AddColumn(helper.NewNumericReportColumn("J", j), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
