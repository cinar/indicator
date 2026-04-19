// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volume"
)

const (
	// DefaultObvStrategyPeriod is the default OBV strategy period.
	DefaultObvStrategyPeriod = 10
)

// ObvStrategy represents the configuration parameters for calculating the On-Balance Volume (OBV) strategy.
// Recommends a Buy action when OBV crosses above its SMA, and recommends a Sell action when OBV crosses below its SMA.
type ObvStrategy struct {
	// Obv is the OBV indicator instance.
	Obv *volume.Obv[float64]

	// Sma is the SMA indicator instance.
	Sma *trend.Sma[float64]
}

// NewObvStrategy function initializes a new OBV strategy instance with the default parameters.
func NewObvStrategy() *ObvStrategy {
	return NewObvStrategyWith(
		DefaultObvStrategyPeriod,
	)
}

// NewObvStrategyWith function initializes a new OBV strategy instance with the given period.
func NewObvStrategyWith(period int) *ObvStrategy {
	return &ObvStrategy{
		Obv: volume.NewObv[float64](),
		Sma: trend.NewSmaWithPeriod[float64](period),
	}
}

// Name function returns the name of the strategy.
func (s *ObvStrategy) Name() string {
	return fmt.Sprintf("OBV Strategy (%d)", s.Sma.Period)
}

// Compute function processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *ObvStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	closings := asset.SnapshotsAsClosings(snapshotsSplice[0])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[1])

	obvValues := s.Obv.Compute(closings, volumes)
	obvSplice := helper.Duplicate(obvValues, 2)

	smaValues := s.Sma.Compute(obvSplice[0])

	// Align OBV with SMA
	obvValuesAligned := helper.Skip(obvSplice[1], s.Sma.IdlePeriod())

	actions := helper.Operate(obvValuesAligned, smaValues, func(obv, sma float64) strategy.Action {
		if obv > sma {
			return strategy.Buy
		}

		if obv < sma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// OBV starts after its idle period (0), but SMA starts after its idle period.
	actions = helper.Shift(actions, s.Sma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report function processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (s *ObvStrategy) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings (for report)
	// snapshots[2] -> closings (for obv)
	// snapshots[3] -> volumes  (for obv)
	// snapshots[4] -> actions / outcomes
	//
	snapshotsSplice := helper.Duplicate(snapshots, 5)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		s.Sma.IdlePeriod(),
	)

	closingsForReport := helper.Duplicate(
		helper.Skip(
			asset.SnapshotsAsClosings(snapshotsSplice[1]),
			s.Sma.IdlePeriod(),
		),
		2,
	)

	closingsForObv := asset.SnapshotsAsClosings(snapshotsSplice[2])
	volumesForObv := asset.SnapshotsAsVolumes(snapshotsSplice[3])

	obvValues := s.Obv.Compute(closingsForObv, volumesForObv)
	obvSplice := helper.Duplicate(obvValues, 2)

	smaValues := s.Sma.Compute(obvSplice[0])
	obvValuesAligned := helper.Skip(obvSplice[1], s.Sma.IdlePeriod())

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshotsSplice[4])
	actions = helper.Skip(actions, s.Sma.IdlePeriod())
	outcomes = helper.Skip(outcomes, s.Sma.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsForReport[0]))

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsForReport[1]), 1)
	report.AddColumn(helper.NewNumericReportColumn("OBV", obvValuesAligned), 1)
	report.AddColumn(helper.NewNumericReportColumn("SMA", smaValues), 1)

	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
