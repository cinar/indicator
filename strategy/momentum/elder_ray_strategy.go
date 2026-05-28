// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

// ElderRayStrategy represents the configuration parameters for calculating the Elder Ray strategy.
// Buy when EMA is rising and Bear Power is negative but rising.
// Sell when EMA is falling and Bull Power is positive but falling.
type ElderRayStrategy struct {
	// ElderRay represents the configuration parameters for calculating the Elder-Ray Index.
	ElderRay *momentum.ElderRay[float64]
}

// NewElderRayStrategy function initializes a new Elder Ray strategy instance with the default parameters.
func NewElderRayStrategy() *ElderRayStrategy {
	return &ElderRayStrategy{
		ElderRay: momentum.NewElderRay[float64](),
	}
}

// Name returns the name of the strategy.
func (*ElderRayStrategy) Name() string {
	return "Elder Ray Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (e *ElderRayStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[2]), 2)

	bullPower, bearPower := e.ElderRay.Compute(highs, lows, closingsSplice[0])

	bullSplice := helper.Duplicate(bullPower, 2)
	bearSplice := helper.Duplicate(bearPower, 2)

	bullChange := helper.Change(bullSplice[0], 1)
	bearChange := helper.Change(bearSplice[0], 1)

	bullCurrent := helper.Skip(bullSplice[1], 1)
	bearCurrent := helper.Skip(bearSplice[1], 1)

	ema := trend.NewEmaWithPeriod[float64](e.ElderRay.Period)
	emaChange := helper.Change(ema.Compute(closingsSplice[1]), 1)

	actions := helper.Operate5(
		bullCurrent, bullChange,
		bearCurrent, bearChange,
		emaChange,
		func(bull, dBull, bear, dBear, dEma float64) strategy.Action {
			// Buy: EMA rising AND Bear Power negative but rising
			if dEma > 0 && bear < 0 && dBear > 0 {
				return strategy.Buy
			}

			// Sell: EMA falling AND Bull Power positive but falling
			if dEma < 0 && bull > 0 && dBull < 0 {
				return strategy.Sell
			}

			return strategy.Hold
		})

	// IdlePeriod + 1 for the Change(1) step
	actions = helper.Shift(actions, e.ElderRay.IdlePeriod()+1, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (e *ElderRayStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs    -|
	// snapshots[2] -> lows     -+-> ElderRay.Compute -> bullPower, bearPower
	// snapshots[3] -> closings -|
	//                 closings -> close
	// snapshots[4] -> actions  -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 5)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[3]), 2)

	bullPower, bearPower := e.ElderRay.Compute(highs, lows, closingsSplice[0])
	bullPower = helper.Shift(bullPower, e.ElderRay.IdlePeriod(), 0)
	bearPower = helper.Shift(bearPower, e.ElderRay.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(e, snapshots[4])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(e.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Bull Power", bullPower), 1)
	report.AddColumn(helper.NewNumericReportColumn("Bear Power", bearPower), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
