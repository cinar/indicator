// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultStochasticOscillatorStrategyBuyAt defines the default K level at which a Buy action is generated.
	DefaultStochasticOscillatorStrategyBuyAt = 20.0

	// DefaultStochasticOscillatorStrategySellAt defines the default K level at which a Sell action is generated.
	DefaultStochasticOscillatorStrategySellAt = 80.0
)

// StochasticOscillatorStrategy represents the configuration parameters for calculating the Stochastic Oscillator
// strategy. When the K line is below the buy threshold, a Buy action is generated. When above the sell threshold,
// a Sell action is generated.
type StochasticOscillatorStrategy struct {
	// StochasticOscillator represents the configuration parameters for calculating the Stochastic Oscillator.
	StochasticOscillator *momentum.StochasticOscillator[float64]

	// BuyAt defines the K level at which a Buy action is generated.
	BuyAt float64

	// SellAt defines the K level at which a Sell action is generated.
	SellAt float64
}

// NewStochasticOscillatorStrategy function initializes a new Stochastic Oscillator strategy instance with
// the default parameters.
func NewStochasticOscillatorStrategy() *StochasticOscillatorStrategy {
	return NewStochasticOscillatorStrategyWith(
		DefaultStochasticOscillatorStrategyBuyAt,
		DefaultStochasticOscillatorStrategySellAt,
	)
}

// NewStochasticOscillatorStrategyWith function initializes a new Stochastic Oscillator strategy instance with
// the given parameters.
func NewStochasticOscillatorStrategyWith(buyAt, sellAt float64) *StochasticOscillatorStrategy {
	return &StochasticOscillatorStrategy{
		StochasticOscillator: momentum.NewStochasticOscillator[float64](),
		BuyAt:                buyAt,
		SellAt:               sellAt,
	}
}

// Name returns the name of the strategy.
func (s *StochasticOscillatorStrategy) Name() string {
	return fmt.Sprintf("Stochastic Oscillator Strategy (%.0f,%.0f)", s.BuyAt, s.SellAt)
}

// ComputeWithContext processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *StochasticOscillatorStrategy) ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.DuplicateWithContext(ctx, snapshots, 3)

	highs := asset.SnapshotsAsHighsWithContext(ctx, snapshotsSplice[0])
	lows := asset.SnapshotsAsLowsWithContext(ctx, snapshotsSplice[1])
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshotsSplice[2])

	k, d := s.StochasticOscillator.ComputeWithContext(ctx, highs, lows, closings)

	var prevK, prevD float64
	var hasPrev bool

	actions := helper.OperateWithContext(ctx, k, d, func(kVal, dVal float64) strategy.Action {
		if !hasPrev {
			prevK = kVal
			prevD = dVal
			hasPrev = true
			return strategy.Hold
		}

		action := strategy.Hold

		if prevK <= prevD && kVal > dVal && kVal < s.BuyAt {
			action = strategy.Buy
		} else if prevK >= prevD && kVal < dVal && kVal > s.SellAt {
			action = strategy.Sell
		}

		prevK = kVal
		prevD = dVal

		return action
	})

	// Stochastic Oscillator starts only after the idle period.
	actions = helper.ShiftWithContext(ctx, actions, s.StochasticOscillator.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (s *StochasticOscillatorStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> highs   -|
	// snapshots[2] -> lows    -+-> StochasticOscillator.Compute -> k, d
	// snapshots[3] -> closings-|
	// snapshots[4] -> closings -> close
	// snapshots[5] -> actions  -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 6)

	dates := asset.SnapshotsAsDates(snapshots[0])
	highs := asset.SnapshotsAsHighs(snapshots[1])
	lows := asset.SnapshotsAsLows(snapshots[2])
	closings := asset.SnapshotsAsClosings(snapshots[3])
	closings2 := asset.SnapshotsAsClosings(snapshots[4])

	k, d := s.StochasticOscillator.Compute(highs, lows, closings)
	k = helper.Shift(k, s.StochasticOscillator.IdlePeriod(), 0)
	d = helper.Shift(d, s.StochasticOscillator.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshots[5])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings2))
	report.AddColumn(helper.NewNumericReportColumn("K", k), 1)
	report.AddColumn(helper.NewNumericReportColumn("D", d), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (s *StochasticOscillatorStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return s.ComputeWithContext(context.Background(), snapshots)
}
