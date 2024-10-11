// Copyright (c) 2021-2024 Onur Cinar.
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

// EnvelopeStrategy represents the configuration parameters for calculating the Envelope strategy. When the closing
// is above the upper band suggests a Sell recommendation, and when the closing is below the lower band suggests a
// buy recommendation.
type EnvelopeStrategy struct {
	// Envelope is the envelope indicator instance.
	Envelope *trend.Envelope[float64]
}

// NewEnvelopeStrategy function initializes a new Envelope strategy with the default parameters.
func NewEnvelopeStrategy() *EnvelopeStrategy {
	return NewEnvelopeStrategyWith(
		trend.NewEnvelopeWithSma[float64](),
	)
}

// NewEnvelopeStrategyWith function initializes a new Envelope strategy with the given Envelope instance.
func NewEnvelopeStrategyWith(envelope *trend.Envelope[float64]) *EnvelopeStrategy {
	return &EnvelopeStrategy{
		Envelope: envelope,
	}
}

// Name returns the name of the strategy.
func (e *EnvelopeStrategy) Name() string {
	return fmt.Sprintf("Envelope Strategy (%s)", e.Envelope.String())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (e *EnvelopeStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots),
		2,
	)

	uppers, middles, lowers := e.Envelope.Compute(closingsSplice[0])
	go helper.Drain(middles)

	actions := helper.Operate3(uppers, lowers, closingsSplice[1], func(upper, lower, closing float64) strategy.Action {
		// When the closing is below the lower band suggests a buy recommendation.
		if closing < lower {
			return strategy.Buy
		}

		// When the closing is above the upper band suggests a Sell recommendation.
		if closing > upper {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Envelope start only after a full period.
	actions = helper.Shift(actions, e.Envelope.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (e *EnvelopeStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> envelope -> upper
	//                                         -> middle
	//                                         -> lower
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshotsSplice := helper.Duplicate(c, 3)

	dates := helper.Skip(
		asset.SnapshotsAsDates(snapshotsSplice[0]),
		e.Envelope.IdlePeriod(),
	)

	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshotsSplice[1]), 2)
	closingsSplice[1] = helper.Skip(closingsSplice[1], t.IdlePeriod())

	tsisSplice := helper.Duplicate(t.Tsi.Compute(closingsSplice[0]), 2)
	tsisSplice[0] = helper.Skip(tsisSplice[0], t.Signal.IdlePeriod())

	signals := t.Signal.Compute(tsisSplice[1])

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshotsSplice[2])
	actions = helper.Skip(actions, t.IdlePeriod())
	outcomes = helper.Skip(outcomes, t.IdlePeriod())

	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closingsSplice[1]))

	report.AddColumn(helper.NewNumericReportColumn("TSI", tsisSplice[0]), 1)
	report.AddColumn(helper.NewNumericReportColumn("Signal", signals), 1)

	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
