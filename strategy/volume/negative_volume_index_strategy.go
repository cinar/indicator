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
	// DefaultNegativeVolumeIndexStrategyEmaPeriod is the default EMA period of 255.
	DefaultNegativeVolumeIndexStrategyEmaPeriod = 255
)

// NegativeVolumeIndexStrategy represents the configuration parameters for calculating the Negative Volume Index
// strategy. Recommends a Buy action when it crosses below its EMA, recommends a Sell action when it crosses
// above its EMA, and recommends a Hold action otherwise.
type NegativeVolumeIndexStrategy struct {
	// NegativeVolumeIndex is the Negative Volume Index indicator instance.
	NegativeVolumeIndex *volume.Nvi[float64]

	// NegativeVolumeIndexEma is the Negative Volume Index EMA instance.
	NegativeVolumeIndexEma *trend.Ema[float64]
}

// NewNegativeVolumeIndexStrategy function initializes a new Negative Volume Index strategy instance with the
// default parameters.
func NewNegativeVolumeIndexStrategy() *NegativeVolumeIndexStrategy {
	return NewNegativeVolumeIndexStrategyWith(
		DefaultNegativeVolumeIndexStrategyEmaPeriod,
	)
}

// NewNegativeVolumeIndexStrategyWith function initializes a new Negative Volume Index strategy instance with the
// given parameters.
func NewNegativeVolumeIndexStrategyWith(emaPeriod int) *NegativeVolumeIndexStrategy {
	return &NegativeVolumeIndexStrategy{
		NegativeVolumeIndex:    volume.NewNvi[float64](),
		NegativeVolumeIndexEma: trend.NewEmaWithPeriod[float64](emaPeriod),
	}
}

// Name returns the name of the strategy.
func (n *NegativeVolumeIndexStrategy) Name() string {
	return fmt.Sprintf("Negative Volume Index Strategy (%d)", n.NegativeVolumeIndexEma.Period)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (n *NegativeVolumeIndexStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	closings := asset.SnapshotsAsClosings(snapshotsSplice[0])
	volumes := asset.SnapshotsAsVolumes(snapshotsSplice[1])

	nvisSplice := helper.Duplicate(
		n.NegativeVolumeIndex.Compute(closings, volumes),
		2,
	)

	nvisSplice[0] = helper.Skip(nvisSplice[0], n.NegativeVolumeIndexEma.IdlePeriod())
	nviEmas := n.NegativeVolumeIndexEma.Compute(nvisSplice[1])

	actions := helper.Operate(nvisSplice[0], nviEmas, func(nvi, nviEma float64) strategy.Action {
		if nvi < nviEma {
			return strategy.Buy
		}

		if nvi > nviEma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Negative Volume Index starts only after a full period.
	actions = helper.Shift(
		actions,
		n.NegativeVolumeIndex.IdlePeriod()+n.NegativeVolumeIndexEma.IdlePeriod(),
		strategy.Hold,
	)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (n *NegativeVolumeIndexStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> negative volume index[0] -> negative volume index
	//                                negative volume index[1] -> negative volume index ema
	// snapshots[2] -> volumes
	// snapshots[3] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 4)

	period := n.NegativeVolumeIndex.IdlePeriod() + n.NegativeVolumeIndexEma.IdlePeriod()

	dates := helper.Skip(asset.SnapshotsAsDates(snapshots[0]), period)

	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots[1]),
		2,
	)
	volumes := asset.SnapshotsAsVolumes(snapshots[2])

	nvisSplice := helper.Duplicate(
		n.NegativeVolumeIndex.Compute(closingsSplice[0], volumes),
		2,
	)

	nvisSplice[0] = helper.Skip(nvisSplice[0], n.NegativeVolumeIndexEma.IdlePeriod())
	nviEmas := n.NegativeVolumeIndexEma.Compute(nvisSplice[1])

	closingsSplice[1] = helper.Skip(closingsSplice[1], period)

	actions, outcomes := strategy.ComputeWithOutcome(n, snapshots[3])
	actions = helper.Skip(actions, period)
	outcomes = helper.Skip(outcomes, period)

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(n.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))
	report.AddColumn(helper.NewNumericReportColumn("NVI", nvisSplice[0]), 1)
	report.AddColumn(helper.NewNumericReportColumn("NVI EMA", nviEmas), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
