// Copyright (c) 2021-2026 Onur Cinar.
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

// ChaikinMoneyFlowStrategy represents the configuration parameters for calculating the Chaikin Money Flow strategy.
// Recommends a Buy action when it crosses above 0, and recommends a Sell action when it crosses below 0.
type ChaikinMoneyFlowStrategy struct {
	// ChaikinMoneyFlow is the Chaikin Money Flow indicator instance.
	ChaikinMoneyFlow *volume.Cmf[float64]
}

// NewChaikinMoneyFlowStrategy function initializes a new Chaikin Money Flow strategy instance with the
// default parameters.
func NewChaikinMoneyFlowStrategy() *ChaikinMoneyFlowStrategy {
	return NewChaikinMoneyFlowStrategyWith(
		volume.DefaultCmfPeriod,
	)
}

// NewChaikinMoneyFlowStrategyWith function initializes a new Chaikin Money Flow strategy instance with the
// given parameters.
func NewChaikinMoneyFlowStrategyWith(period int) *ChaikinMoneyFlowStrategy {
	return &ChaikinMoneyFlowStrategy{
		ChaikinMoneyFlow: volume.NewCmfWithPeriod[float64](period),
	}
}

// Name function returns the name of the strategy.
func (c *ChaikinMoneyFlowStrategy) Name() string {
	return fmt.Sprintf("Chaikin Money Flow Strategy (%d)", c.ChaikinMoneyFlow.IdlePeriod()+1)
}

// Compute function processes the provided asset snapshots and generates a stream of actionable recommendations.
func (c *ChaikinMoneyFlowStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 4)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[0])
	lows := asset.SnapshotsAsLows(snapshotsSplice[1])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[2])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[3])

	cmfs := c.ChaikinMoneyFlow.Compute(highs, lows, closings, volumes)

	actions := helper.Map(cmfs, func(cmf float64) strategy.Action {
		if cmf > 0 {
			return strategy.Buy
		}

		if cmf < 0 {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Chaikin Money Flow starts only after a full period.
	actions = helper.Shift(actions, c.ChaikinMoneyFlow.IdlePeriod(), strategy.Hold)

	return actions
}

// Report function processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (c *ChaikinMoneyFlowStrategy) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs       |
	// snapshots[2] -> lows        |
	// snapshots[3] -> closings[0] -> closings
	//                 closings[1] -> chaikin money flow
	// snapshots[4] -> volumes
	// snapshots[5] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(snapshots, 6)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		c.ChaikinMoneyFlow.IdlePeriod(),
	)

	highs := asset.SnapshotsAsHighs(snapshotsSplice[1])
	lows := asset.SnapshotsAsLows(snapshotsSplice[2])
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshotsSplice[3]),
		2,
	)
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[4])

	cmfs := c.ChaikinMoneyFlow.Compute(highs, lows, closingsSplice[0], volumes)
	closingsSplice[1] = helper.Skip(closingsSplice[1], c.ChaikinMoneyFlow.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(c, snapshotsSplice[5])
	actions = helper.Skip(actions, c.ChaikinMoneyFlow.IdlePeriod())
	outcomes = helper.Skip(outcomes, c.ChaikinMoneyFlow.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(c.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("Chaikin Money Flow", cmfs), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
