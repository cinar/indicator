// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

// AwesomeOscillatorStrategy represents the configuration parameters for calculating the Awesome Oscillator strategy.
type AwesomeOscillatorStrategy struct {
	strategy.Strategy

	// AwesomeOscillator represents the configuration parameters for calculating the Awesome Oscillator.
	AwesomeOscillator *momentum.AwesomeOscillator[float64]
}

// NewAwesomeOscillatorStrategy function initializes a new Awesome Oscillator strategy with the default parameters.
func NewAwesomeOscillatorStrategy() *AwesomeOscillatorStrategy {
	return &AwesomeOscillatorStrategy{
		AwesomeOscillator: momentum.NewAwesomeOscillator[float64](),
	}
}

// Name returns the name of the strategy.
func (*AwesomeOscillatorStrategy) Name() string {
	return "Awesome Oscillator Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (a *AwesomeOscillatorStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])

	ao := a.AwesomeOscillator.Compute(highs, lows)

	actions := helper.Map(ao, func(value float64) strategy.Action {
		if value < 0 {
			return strategy.Sell
		}

		if value > 0 {
			return strategy.Buy
		}

		return strategy.Hold
	})

	// Awesome Oscillator starts only after the idle period.
	actions = helper.Shift(actions, a.AwesomeOscillator.IdlePeriod(), strategy.Hold)

	actions = strategy.NormalizeActions(actions)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (a *AwesomeOscillatorStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute     -> actions -> annotations
	// snapshots[2] -> closings[0] -> close
	// snapshots[3] -> highs -|
	// snapshots[4] -> lows  -> AwesomeOscillator.Compute -> ao
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[2])
	highs := asset.SnapshotsAsHighs(snapshots[3])
	lows := asset.SnapshotsAsLows(snapshots[4])

	ao := helper.Shift(a.AwesomeOscillator.Compute(highs, lows), a.AwesomeOscillator.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(a, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(a.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("AO", ao), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
