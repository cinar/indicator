// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"context"

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

// AlligatorStrategy demonstrates how to compose three Smoothed Moving Averages (SMMAs)
// (jaw, teeth, lip) into an illustrative multi-moving-average trend-following strategy.
type AlligatorStrategy struct {
	// Jaw represents the slowest moving aveage.
	Jaw *trend.Smma[float64]

	// Teeth represents the medium moving average.
	Teeth *trend.Smma[float64]

	// Lip represents the fastest moving average.
	Lip *trend.Smma[float64]
}

// NewAlligatorStrategy initializes an example AlligatorStrategy instance with default parameters.
func NewAlligatorStrategy() *AlligatorStrategy {
	return NewAlligatorStrategyWith(
		DefaultAlligatorStrategyJawPeriod,
		DefaultAlligatorStrategyTeethPeriod,
		DefaultAlligatorStrategyLipPeriod,
	)
}

// NewAlligatorStrategyWith initializes an example AlligatorStrategyWith instance with default parameters.
func NewAlligatorStrategyWith(jawPeriod, teethPeriod, lipPeriod int) *AlligatorStrategy {
	return &AlligatorStrategy{
		Jaw:   trend.NewSmmaWithPeriod[float64](jawPeriod),
		Teeth: trend.NewSmmaWithPeriod[float64](teethPeriod),
		Lip:   trend.NewSmmaWithPeriod[float64](lipPeriod),
	}
}

// Name returns the name of the example strategy.
func (a *AlligatorStrategy) Name() string {
	return fmt.Sprintf("Alligator Strategy (%d,%d,%d)",
		a.Jaw.Period,
		a.Teeth.Period,
		a.Lip.Period,
	)
}

// ComputeWithContext processes the provided asset snapshots and generates an illustrative stream of actions.
func (a *AlligatorStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.DuplicateWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshots), 3)

	jaws := a.Jaw.ComputeWithContext(ctx, closingsSplice[0])
	teeths := a.Teeth.ComputeWithContext(ctx, closingsSplice[1])
	lips := a.Lip.ComputeWithContext(ctx, closingsSplice[2])

	commonPeriod := helper.CommonPeriod(a.Jaw.Period, a.Teeth.Period, a.Lip.Period)
	jaws = helper.SyncPeriod(commonPeriod, a.Jaw.Period, jaws)
	teeths = helper.SyncPeriod(commonPeriod, a.Teeth.Period, teeths)
	lips = helper.SyncPeriod(commonPeriod, a.Lip.Period, lips)

	actions := helper.Operate3WithContext(ctx, jaws, teeths, lips, func(jaw, teeth, lip float64) strategy.Action {
		if lip > teeth && lip > jaw {
			return strategy.Buy
		}

		if lip < teeth && lip < jaw {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Alligator strategy starts only after a full period.
	actions = helper.ShiftWithContext(ctx, actions, commonPeriod, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates an
// illustrative report annotated with example actions.
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

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (a *AlligatorStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return a.ComputeWithContext(context.Background(), snapshots)
}
