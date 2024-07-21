// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultStochasticRsiStrategyBuyAt defines the default level at which a Buy action is generated.
	DefaultStochasticRsiStrategyBuyAt = 0.8

	// DefaultStochasticRsiStrategySellAt defines the default level at which a Sell action is generated.
	DefaultStochasticRsiStrategySellAt = 0.2
)

// StochasticRsiStrategy represents the configuration parameters for calculating the Stochastic RSI strategy.
type StochasticRsiStrategy struct {
	// StochasticRsi represents the configuration parameters for calculating the Stochastic RSI.
	StochasticRsi *momentum.StochasticRsi[float64]

	// BuyAt defines the level at which a Buy action is generated.
	BuyAt float64

	// SellAt defines the level at which a Sell action is generated.
	SellAt float64
}

// NewStochasticRsiStrategy function initializes a new Stochastic RSI strategy instance with the default parameters.
func NewStochasticRsiStrategy() *StochasticRsiStrategy {
	return NewStochasticRsiStrategyWith(
		DefaultStochasticRsiStrategyBuyAt,
		DefaultStochasticRsiStrategySellAt,
	)
}

// NewStochasticRsiStrategyWith function initializes a new Stochastic RSI strategy instance with the given parameters.
func NewStochasticRsiStrategyWith(buyAt, sellAt float64) *StochasticRsiStrategy {
	return &StochasticRsiStrategy{
		StochasticRsi: momentum.NewStochasticRsi[float64](),
		BuyAt:         buyAt,
		SellAt:        sellAt,
	}
}

// Name returns the name of the strategy.
func (s *StochasticRsiStrategy) Name() string {
	return fmt.Sprintf("Stochastic RSI Strategy (%.1f,%.1f)", s.BuyAt, s.SellAt)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *StochasticRsiStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := asset.SnapshotsAsClosings(snapshots)

	stochasticRsi := s.StochasticRsi.Compute(closings)

	actions := helper.Map(stochasticRsi, func(value float64) strategy.Action {
		if value <= s.BuyAt {
			return strategy.Buy
		}

		if value >= s.SellAt {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// Stochastic RSI starts only after the idle period.
	actions = helper.Shift(actions, s.StochasticRsi.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (s *StochasticRsiStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute     -> actions -> annotations
	// snapshots[2] -> closings[0] -> close
	//              -> closings[1] -> StochasticRsi.Compute -> stochasticRsi
	//
	snapshots := helper.Duplicate(c, 3)

	dates := helper.Skip(asset.SnapshotsAsDates(snapshots[0]), s.StochasticRsi.IdlePeriod())

	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 2)
	closings[0] = helper.Skip(closings[0], s.StochasticRsi.IdlePeriod())

	stochasticRsi := s.StochasticRsi.Compute(closings[1])

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshots[1])
	annotations := helper.Skip(strategy.ActionsToAnnotations(actions), s.StochasticRsi.IdlePeriod())
	outcomes = helper.Skip(helper.MultiplyBy(outcomes, 100), s.StochasticRsi.IdlePeriod())

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("Stochastic RSI", stochasticRsi), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
