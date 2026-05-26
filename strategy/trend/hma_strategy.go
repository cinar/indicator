// Copyright (c) 2021-2026 Onur Cinar.
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
	// DefaultHmaStrategyPeriod is the default period for the HMA strategy.
	DefaultHmaStrategyPeriod = 9
)

// HmaStrategy represents the configuration parameters for calculating the HMA strategy. A closing price crossing
// above the HMA suggests a bullish trend, while crossing below the HMA indicates a bearish trend.
type HmaStrategy struct {
	// Hma represents the configuration parameters for calculating the Hull Moving Average.
	Hma *trend.Hma[float64]
}

// NewHmaStrategy function initializes a new HMA strategy instance with the default parameters.
func NewHmaStrategy() *HmaStrategy {
	return NewHmaStrategyWith(DefaultHmaStrategyPeriod)
}

// NewHmaStrategyWith function initializes a new HMA strategy instance with the given period.
func NewHmaStrategyWith(period int) *HmaStrategy {
	return &HmaStrategy{
		Hma: trend.NewHmaWithPeriod[float64](period),
	}
}

// Name returns the name of the strategy.
func (h *HmaStrategy) Name() string {
	return h.Hma.String() + " Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (h *HmaStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshots), 2)
	closingsSplice[1] = helper.Skip(closingsSplice[1], h.Hma.IdlePeriod())

	hmas := h.Hma.Compute(closingsSplice[0])

	actions := helper.Operate(hmas, closingsSplice[1], func(hma, closing float64) strategy.Action {
		if closing > hma {
			return strategy.Buy
		}

		if closing < hma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// HMA starts only after a full period.
	actions = helper.Shift(actions, h.Hma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (h *HmaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> hma
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshotsSplice[0])
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[1]), 2)

	hmas := h.Hma.Compute(closingsSplice[0])
	hmas = helper.Shift(hmas, h.Hma.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(h, snapshotsSplice[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(h.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("HMA", hmas))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
