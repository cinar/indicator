// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/volume"
)

// WeightedAveragePriceStrategy represents the configuration parameters for calculating the Weighted
// Average Price strategy. Recommends a Buy action when the closing crosses below the VWAP, recommends a Sell
// action when the closing crosses above the VWAP, and recommends a Hold action otherwise.
type WeightedAveragePriceStrategy struct {
	// WeightedAveragePrice is the Weighted Average Price indicator instance.
	WeightedAveragePrice *volume.Vwap[float64]
}

// NewWeightedAveragePriceStrategy function initializes a new Weighted Average Price strategy
// instance with the default parameters.
func NewWeightedAveragePriceStrategy() *WeightedAveragePriceStrategy {
	return NewWeightedAveragePriceStrategyWith(
		volume.DefaultVwapPeriod,
	)
}

// NewWeightedAveragePriceStrategyWith function initializes a new Weighted Average Price strategy
// instance with the given parameters.
func NewWeightedAveragePriceStrategyWith(period int) *WeightedAveragePriceStrategy {
	return &WeightedAveragePriceStrategy{
		WeightedAveragePrice: volume.NewVwapWithPeriod[float64](period),
	}
}

// Name returns the name of the strategy.
func (v *WeightedAveragePriceStrategy) Name() string {
	return fmt.Sprintf("Weighted Average Price Strategy (%d)", v.WeightedAveragePrice.IdlePeriod()+1)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (v *WeightedAveragePriceStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshotsSplice[0]),
		2,
	)

	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[1])

	vwaps := v.WeightedAveragePrice.Compute(closingsSplice[1], volumes)
	closingsSplice[0] = helper.Skip(closingsSplice[0], v.WeightedAveragePrice.IdlePeriod())

	actions := helper.Operate(closingsSplice[0], vwaps, func(closing, vwap float64) strategy.Action {
		if vwap > closing {
			return strategy.Buy
		}

		if vwap < closing {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Weighted Average Price starts only after a full period.
	actions = helper.Shift(actions, v.WeightedAveragePrice.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (v *WeightedAveragePriceStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> vwap
	// snapshots[2] -> volumes
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	dates := helper.Skip(asset.SnapshotsAsDates(snapshots[0]), v.WeightedAveragePrice.IdlePeriod())

	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots[1]),
		2,
	)
	volumes := asset.SnapshotsAsVolumes(snapshots[2])

	vwaps := v.WeightedAveragePrice.Compute(closingsSplice[0], volumes)

	closingsSplice[1] = helper.Skip(closingsSplice[1], v.WeightedAveragePrice.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(v, snapshots[3])
	actions = helper.Skip(actions, v.WeightedAveragePrice.IdlePeriod())
	outcomes = helper.Skip(outcomes, v.WeightedAveragePrice.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(v.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("VWAP", vwaps))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
