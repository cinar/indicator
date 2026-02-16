// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultAlligatorStrategyJawPeriod is the default jaw period of 13.
	DefaultAlligatorStrategyJawPeriod = 13

	// DefaultAlligatorStrategyTeethPeriod is the default teeth period of 8.
	DefaultAlligatorStrategyTeethPeriod = 8

	// DefaultAlligatorStrategyLipPeriod is the default lip period of 5.
	DefaultAlligatorStrategyLipPeriod = 5
)

// AlligatorStrategy represents the configuration parameters for calculating the
// Alligator strategy. It is a technical indicator to help identify the presence
// and the direction of the trend. It uses three Smooted Moving Averges (SMMAs).
type AlligatorStrategy struct {
	// Jaw represents the slowest moving aveage.
	Jaw *trend.Smma[float64]

	// Teeth represents the medium moving average.
	Teeth *trend.Smma[float64]

	// Lip represents the fastest moving average.
	Lip *trend.Smma[float64]
}

// NewAlligatorStrategy function initializes a new Alligator strategy instance.
func NewAlligatorStrategy() *AlligatorStrategy {
	return NewAlligatorStrategyWith(
		DefaultAlligatorStrategyJawPeriod,
		DefaultAlligatorStrategyTeethPeriod,
		DefaultAlligatorStrategyLipPeriod,
	)
}

// NewAlligatorStrategyWith function initializes a new Alligator strategy instance with the given parameters.
func NewAlligatorStrategyWith(jawPeriod, teethPeriod, lipPeriod int) *AlligatorStrategy {
	return &AlligatorStrategy{
		Jaw:   trend.NewSmmaWithPeriod[float64](jawPeriod),
		Teeth: trend.NewSmmaWithPeriod[float64](teethPeriod),
		Lip:   trend.NewSmmaWithPeriod[float64](lipPeriod),
	}
}

// Name returns the name of the strategy.
func (a *AlligatorStrategy) Name() string {
	return fmt.Sprintf("Alligator Strategy (%d,%d,%d)",
		a.Jaw.Period,
		a.Teeth.Period,
		a.Lip.Period,
	)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (a *AlligatorStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshots), 3)

	jaws := a.Jaw.Compute(closingsSplice[0])
	teeths := a.Teeth.Compute(closingsSplice[1])
	lips := a.Lip.Compute(closingsSplice[2])

	commonPeriod := helper.CommonPeriod(a.Jaw.Period, a.Teeth.Period, a.Lip.Period)
	jaws = helper.SyncPeriod(commonPeriod, a.Jaw.Period, jaws)
	teeths = helper.SyncPeriod(commonPeriod, a.Teeth.Period, teeths)
	lips = helper.SyncPeriod(commonPeriod, a.Lip.Period, lips)

	actions := helper.Operate3(jaws, teeths, lips, func(jaw, teeth, lip float64) strategy.Action {
		if lip > teeth && lip > jaw {
			return strategy.Buy
		}

		if lip < teeth && lip < jaw {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Alligator strategy starts only after a full period.
	actions = helper.Shift(actions, commonPeriod, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (a *AlligatorStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> jaw
	//                 closings[2] -> teeth
	//                 closings[3] -> lip
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 4)

	jaws := a.Jaw.Compute(closingsSplice[1])
	teeths := a.Teeth.Compute(closingsSplice[2])
	lips := a.Lip.Compute(closingsSplice[3])

	actions, outcomes := strategy.ComputeWithOutcome(a, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	commonPeriod := helper.CommonPeriod(a.Jaw.Period, a.Teeth.Period, a.Lip.Period)
	dates = helper.SyncPeriod(commonPeriod, 0, dates)
	closingsSplice[0] = helper.Skip(closingsSplice[0], commonPeriod)
	jaws = helper.SyncPeriod(commonPeriod, a.Jaw.Period, jaws)
	teeths = helper.SyncPeriod(commonPeriod, a.Teeth.Period, teeths)
	lips = helper.SyncPeriod(commonPeriod, a.Lip.Period, lips)
	annotations = helper.Skip(annotations, commonPeriod)
	outcomes = helper.Skip(outcomes, commonPeriod)

	report := helper.NewReport(a.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[0]))
	report.AddColumn(helper.NewNumericReportColumn("Jaw", jaws), 1)
	report.AddColumn(helper.NewNumericReportColumn("Teeth", teeths), 1)
	report.AddColumn(helper.NewNumericReportColumn("Lip", lips), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
