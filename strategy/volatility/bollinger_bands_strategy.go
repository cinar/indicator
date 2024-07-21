// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volatility"
)

// BollingerBandsStrategy represents the configuration parameters for calculating the Bollinger Bands strategy.
// A closing value crossing above the upper band suggets a Buy signal, while crossing below the lower band
// indivates a Sell signal.
type BollingerBandsStrategy struct {
	strategy.Strategy

	// BollingerBands represents the configuration parameters for calculating the Bollinger Bands.
	BollingerBands *volatility.BollingerBands[float64]
}

// NewBollingerBandsStrategy function initializes a new Bollinger Bands strategy instance.
func NewBollingerBandsStrategy() *BollingerBandsStrategy {
	return &BollingerBandsStrategy{
		BollingerBands: volatility.NewBollingerBands[float64](),
	}
}

// Name returns the name of the strategy.
func (*BollingerBandsStrategy) Name() string {
	return "Bollinger Bands Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (b *BollingerBandsStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots),
		2,
	)

	uppers, middles, lowers := b.BollingerBands.Compute(closings[0])
	go helper.Drain(middles)

	closings[1] = helper.Skip(closings[1], b.BollingerBands.IdlePeriod())

	actions := helper.Operate3(uppers, lowers, closings[1], func(upper, lower, closing float64) strategy.Action {
		if closing > upper {
			return strategy.Buy
		}

		if lower > closing {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Bollinger Bands starts only after a full period.
	actions = helper.Shift(actions, b.BollingerBands.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (b *BollingerBandsStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> upper
	//                             -> middle
	//                             -> lower
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 2)

	uppers, middles, lowers := b.BollingerBands.Compute(closings[0])
	uppers = helper.Shift(uppers, b.BollingerBands.IdlePeriod(), 0)
	middles = helper.Shift(middles, b.BollingerBands.IdlePeriod(), 0)
	lowers = helper.Shift(lowers, b.BollingerBands.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(b, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(b.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("Upper", uppers))
	report.AddColumn(helper.NewNumericReportColumn("Middle", middles))
	report.AddColumn(helper.NewNumericReportColumn("Lower", lowers))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
