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

// IchimokuCloudStrategy represents the configuration parameters for calculating the Ichimoku Cloud strategy.
type IchimokuCloudStrategy struct {
	// IchimokuCloud represents the configuration parameters for calculating the Ichimoku Cloud.
	IchimokuCloud *momentum.IchimokuCloud[float64]
}

// NewIchimokuCloudStrategy function initializes a new Ichimoku Cloud strategy with the default parameters.
func NewIchimokuCloudStrategy() *IchimokuCloudStrategy {
	return &IchimokuCloudStrategy{
		IchimokuCloud: momentum.NewIchimokuCloud[float64](),
	}
}

// Name returns the name of the strategy.
func (*IchimokuCloudStrategy) Name() string {
	return "Ichimoku Cloud Strategy"
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (i *IchimokuCloudStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 3)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])

	closingsSplice := helper.Duplicate(closings, 2)

	cl, bl, lsa, lsb, ll := i.IchimokuCloud.Compute(highs, lows, closingsSplice[0])

	// Lagging line is not used in the core logic, drain it to prevent blocking
	go helper.Drain(ll)

	actions := helper.Operate5(
		helper.Skip(closingsSplice[1], i.IchimokuCloud.IdlePeriod()),
		cl,
		bl,
		lsa,
		lsb,
		func(c, conversion, base, spanA, spanB float64) strategy.Action {
			if c > spanA && c > spanB && conversion > base && spanA > spanB {
				return strategy.Buy
			}

			if c < spanA && c < spanB && conversion < base && spanA < spanB {
				return strategy.Sell
			}

			return strategy.Hold
		},
	)

	// Shift the actions to account for the idle period
	return helper.Shift(actions, i.IchimokuCloud.IdlePeriod(), strategy.Hold)
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (i *IchimokuCloudStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[2])
	highs := asset.SnapshotsAsHighs(snapshots[3])
	lows := asset.SnapshotsAsLows(snapshots[4])
	closingsForCloud := asset.SnapshotsAsClosings(snapshots[5])

	cl, bl, lsa, lsb, ll := i.IchimokuCloud.Compute(highs, lows, closingsForCloud)

	// Lagging line is not used in the report right now, drain it.
	go helper.Drain(ll)

	clShifted := helper.Shift(cl, i.IchimokuCloud.IdlePeriod(), 0)
	blShifted := helper.Shift(bl, i.IchimokuCloud.IdlePeriod(), 0)
	lsaShifted := helper.Shift(lsa, i.IchimokuCloud.IdlePeriod(), 0)
	lsbShifted := helper.Shift(lsb, i.IchimokuCloud.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(i, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(i.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewNumericReportColumn("Conversion Line", clShifted))
	report.AddColumn(helper.NewNumericReportColumn("Base Line", blShifted))
	report.AddColumn(helper.NewNumericReportColumn("Leading Span A", lsaShifted))
	report.AddColumn(helper.NewNumericReportColumn("Leading Span B", lsbShifted))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
