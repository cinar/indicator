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
	// DefaultVwmaStrategyPeriod is the default VWMA period.
	DefaultVwmaStrategyPeriod = 20
)

// VwmaStrategy demonstrates how to compose Volume Weighted Moving Average (VWMA)
// and Simple Moving Average (SMA) into an illustrative moving average crossover strategy.
type VwmaStrategy struct {
	// VWMA indicator.
	Vwma *trend.Vwma[float64]

	// SMA indicator.
	Sma *trend.Sma[float64]
}

// NewVwmaStrategy initializes an example VwmaStrategy instance with default parameters.
func NewVwmaStrategy() *VwmaStrategy {
	v := &VwmaStrategy{
		Vwma: trend.NewVwma[float64](),
		Sma:  trend.NewSma[float64](),
	}

	v.Vwma.Period = DefaultVwmaStrategyPeriod
	v.Sma.Period = DefaultVwmaStrategyPeriod

	return v
}

// Name returns the name of the example strategy.
func (*VwmaStrategy) Name() string {
	return "VWMA Strategy"
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (v *VwmaStrategy) ComputeWithContext(ctx context.Context, c <-chan *asset.Snapshot) <-chan strategy.Action {
	smas, vwmas := v.calculateSmaAndVwma(c)

	actions := helper.OperateWithContext(ctx, smas, vwmas, func(sma, vwma float64) strategy.Action {
		if vwma > sma {
			return strategy.Buy
		}

		if sma > vwma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// VWMA starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, v.Vwma.Period-1, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
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
	volume := asset.SnapshotsAsVolumes(snapshots[1])

	smas := v.Sma.Compute(closings[0])
	vwmas := v.Vwma.Compute(closings[1], volume)

	return smas, vwmas
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (v *VwmaStrategy) Compute(c <-chan *asset.Snapshot) <-chan strategy.Action {
	return v.ComputeWithContext(context.Background(), c)
}
