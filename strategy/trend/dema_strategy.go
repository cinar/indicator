// Copyright (c) 2021-2024 Onur Cinar.
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
	// DefaultDemaStrategyPeriod1 is the first DEMA period.
	DefaultDemaStrategyPeriod1 = 5

	// DefaultDemaStrategyPeriod2 is the second DEMA period.
	DefaultDemaStrategyPeriod2 = 35
)

// DemaStrategy represents the configuration parameters for calculating the DEMA strategy.
// A bullish cross occurs when DEMA with 5 days period moves above DEMA with 35 days period.
// A bearish cross occurs when DEMA with 35 days period moves above DEMA With 5 days period.
type DemaStrategy struct {
	strategy.Strategy

	// Dema1 represents the configuration parameters for
	// calculating the first DEMA.
	Dema1 *trend.Dema[float64]

	// Dema2 represents the configuration parameters for
	// calculating the second DEMA.
	Dema2 *trend.Dema[float64]
}

// NewDemaStrategy function initializes a new DEMA strategy instance
// with the default parameters.
func NewDemaStrategy() *DemaStrategy {
	dema1 := trend.NewDema[float64]()
	dema1.Ema1.Period = DefaultDemaStrategyPeriod1
	dema1.Ema2.Period = DefaultDemaStrategyPeriod1

	dema2 := trend.NewDema[float64]()
	dema2.Ema1.Period = DefaultDemaStrategyPeriod2
	dema2.Ema2.Period = DefaultDemaStrategyPeriod2

	return &DemaStrategy{
		Dema1: dema1,
		Dema2: dema2,
	}
}

// Name returns the name of the strategy.
func (*DemaStrategy) Name() string {
	return "DEMA Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (d *DemaStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := helper.Duplicate(asset.SnapshotsAsClosings(c), 2)

	demas1 := d.Dema1.Compute(closings[0])
	demas1 = helper.Shift(demas1, d.Dema1.IdlePeriod(), 0)

	demas2 := d.Dema2.Compute(closings[1])
	demas2 = helper.Shift(demas2, d.Dema2.IdlePeriod(), 0)

	actions := helper.Operate(demas1, demas2, func(dema1, dema2 float64) strategy.Action {
		if dema1 > dema2 {
			return strategy.Buy
		}

		if dema2 > dema1 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	actions = strategy.NormalizeActions(actions)

	// DEMA starts only after the a full periods for each EMA used.
	actions = helper.Skip(actions, d.Dema2.IdlePeriod())
	actions = helper.Shift(actions, d.Dema2.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (d *DemaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> demas1
	//                 closings[1] -> demas2
	//                 closings[2] -> closings
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 3)

	demas1 := d.Dema1.Compute(closings[0])
	demas1 = helper.Shift(demas1, d.Dema1.IdlePeriod(), 0)

	demas2 := d.Dema2.Compute(closings[1])
	demas2 = helper.Shift(demas2, d.Dema2.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(d, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(d.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[2]))
	report.AddColumn(helper.NewNumericReportColumn(fmt.Sprintf("Dema %d-day", d.Dema1.Ema1.Period), demas1), 1)
	report.AddColumn(helper.NewNumericReportColumn(fmt.Sprintf("Dema %d-day", d.Dema2.Ema1.Period), demas2), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
