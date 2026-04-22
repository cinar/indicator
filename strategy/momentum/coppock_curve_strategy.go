// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// CoppockCurveStrategy represents the configuration parameters for calculating the Coppock Curve strategy.
// A positive Coppock Curve value suggests a Buy signal, while a negative value suggests a Sell signal.
type CoppockCurveStrategy struct {
	// CoppockCurve represents the configuration parameters for calculating the Coppock Curve.
	CoppockCurve *momentum.CoppockCurve[float64]
}

// NewCoppockCurveStrategy function initializes a new Coppock Curve strategy instance with the default parameters.
func NewCoppockCurveStrategy() *CoppockCurveStrategy {
	return &CoppockCurveStrategy{
		CoppockCurve: momentum.NewCoppockCurve[float64](),
	}
}

// Name returns the name of the strategy.
func (*CoppockCurveStrategy) Name() string {
	return "Coppock Curve Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (c *CoppockCurveStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosings(snapshots)

	coppock := c.CoppockCurve.Compute(closings)

	actions := helper.Map(coppock, func(value float64) strategy.Action {
		if value > 0 {
			return strategy.Buy
		}

		if value < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	actions = helper.Shift(actions, c.CoppockCurve.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (c *CoppockCurveStrategy) Report(cr <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute           -> actions -> annotations
	// snapshots[2] -> closings[0]       -> close
	//                 closings[1]       -> CoppockCurve.Compute -> coppock
	//
	snapshots := helper.Duplicate(cr, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])

	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 2)

	coppock := helper.Shift(c.CoppockCurve.Compute(closings[0]), c.CoppockCurve.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(c, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(c.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("Coppock Curve", coppock), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
