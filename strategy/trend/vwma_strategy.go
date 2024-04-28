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

const (
	// DefaultVwmaStrategyPeriod is the default VWMA period.
	DefaultVwmaStrategyPeriod = 20
)

// VwmaStrategy represents the configuration parameters for calculating the VWMA strategy.
// The VwmaStrategy function uses SMA and VWMA indicators to provide a BUY action when
// VWMA is above SMA, and a SELL signal when VWMA is below SMA, a HOLD otherwse.
type VwmaStrategy struct {
	strategy.Strategy

	// VWMA indicator.
	Vwma *trend.Vwma[float64]

	// SMA indicator.
	Sma *trend.Sma[float64]
}

// NewVwmaStrategy function initializes a new VWMA strategy instance with the default parameters.
func NewVwmaStrategy() *VwmaStrategy {
	v := &VwmaStrategy{
		Vwma: trend.NewVwma[float64](),
		Sma:  trend.NewSma[float64](),
	}

	v.Vwma.Period = DefaultVwmaStrategyPeriod
	v.Sma.Period = DefaultVwmaStrategyPeriod

	return v
}

// Name returns the name of the strategy.
func (*VwmaStrategy) Name() string {
	return "VWMA Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (v *VwmaStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	smas, vwmas := v.calculateSmaAndVwma(c)

	actions := helper.Operate(smas, vwmas, func(sma, vwma float64) strategy.Action {
		if vwma > sma {
			return strategy.Buy
		}

		if sma > vwma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	actions = strategy.NormalizeActions(actions)

	// VWMA starts only after the a full period.
	actions = helper.Shift(actions, v.Vwma.Period-1, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (v *VwmaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings
	// snapshots[2] -> sma
	//                 vwma
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	smas, vwmas := v.calculateSmaAndVwma(snapshots[2])
	smas = helper.Shift(smas, v.Vwma.Period-1, 0)
	vwmas = helper.Shift(vwmas, v.Vwma.Period-1, 0)

	actions, outcomes := strategy.ComputeWithOutcome(v, snapshots[3])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(v.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("SMA", smas), 1)
	report.AddColumn(helper.NewNumericReportColumn("VWMA", vwmas), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// calculateSmaAndVwma calculates the SMA and VWMA using the given channel of snapshots.
func (v *VwmaStrategy) calculateSmaAndVwma(c <-chan *asset.Snapshot) (<-chan float64, <-chan float64) {
	snapshots := helper.Duplicate(c, 2)

	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[0]), 2)
	volume := helper.Map(snapshots[1], func(s *asset.Snapshot) float64 { return float64(s.Volume) })

	smas := v.Sma.Compute(closings[0])
	vwmas := v.Vwma.Compute(closings[1], volume)

	return smas, vwmas
}
