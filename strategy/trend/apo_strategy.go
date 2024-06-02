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

// ApoStrategy represents the configuration parameters for calculating the APO strategy.
// An APO value crossing above zero suggests a bullish trend, while crossing below zero
// indicates a bearish trend. Positive APO values signify an upward trend, while
// negative values signify a downward trend.
type ApoStrategy struct {
	strategy.Strategy

	// Apo represents the configuration parameters for calculating the
	// Absolute Price Oscillator (APO).
	Apo *trend.Apo[float64]
}

// NewApoStrategy function initializes a new APO strategy instance with the default parameters.
func NewApoStrategy() *ApoStrategy {
	return &ApoStrategy{
		Apo: trend.NewApo[float64](),
	}
}

// Name returns the name of the strategy.
func (*ApoStrategy) Name() string {
	return "Apo Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (a *ApoStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosings(snapshots)

	apo := a.Apo.Compute(closings)
	apo = helper.Buffered(apo, 2)

	inputs := helper.Duplicate(apo, 2)

	// Skip the first value
	inputs[1] = helper.Skip(inputs[1], 1)

	actions := helper.Operate(inputs[0], inputs[1], func(b, c float64) strategy.Action {
		// An APO value crossing above zero suggests a bullish trend.
		if c >= 0 && b < 0 {
			return strategy.Buy
		}

		// An APO value crossing below zero indicates a bearish trend.
		if c <= 0 && b > 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// APO starts only after the slow period.
	actions = helper.Shift(actions, a.Apo.SlowPeriod, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (a *ApoStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute     -> actions -> annotations
	// snapshots[2] -> closings[0] -> close
	//              -> closings[1] -> Apo.Compute -> apo
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 2)
	apo := helper.Shift(a.Apo.Compute(closings[1]), a.Apo.SlowPeriod, 0)

	actions, outcomes := strategy.ComputeWithOutcome(a, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(a.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("APO", apo), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
