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
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultTripleRsiStrategyPeriod defines the default period for the RSI.
	DefaultTripleRsiStrategyPeriod = 5

	// DefaultTripleRsiStrategyMovingAveragePeriod defines the default period for the SMA.
	DefaultTripleRsiStrategyMovingAveragePeriod = 200

	// DefaultTripleRsiStrategyDownDays defines the default number of down days for the RSI.
	DefaultTripleRsiStrategyDownDays = 3

	// DefaultTripleRsiStrategyBuySignalAt defines the default RSI level at which a Buy signal is confirmed.
	DefaultTripleRsiStrategyBuySignalAt = 60

	// DefaultTripleRsiStrategyBuyAt defines the default RSI level at which a Buy action is generated.
	DefaultTripleRsiStrategyBuyAt = 30

	// DefaultTripleRsiStrategySellAt defines the default RSI level at which a Sell action is generated.
	DefaultTripleRsiStrategySellAt = 50
)

// TripleRsiStrategy represents the configuration parameters for calculating the Triple RSI strategy.
// It assumes that the moving average period is longer than the RSI period.
//
// Recommend Buy:
// - The 5-period RSI is below 30.
// - The 5-period RSI reading is down for the 3rd period in a row.
// - The 5-period RSI reading was below 60 three trading periods ago.
// - The close is higher than the 200-period moving average.
//
// Recommend Sell:
// - Sell at the close when the 5-period RSI crosses above 50.
//
// Based on [Triple RSI Trading Strategy: Enhance Your Win Rate to 90% — Advanced Insights](https://tradingstrategy.medium.com/triple-rsi-trading-strategy-enhance-your-win-rate-to-90-advanced-insights-6143059ce41d).
type TripleRsiStrategy struct {
	// Rsi represents the configuration parameters for calculating the Relative Strength Index (RSI).
	Rsi *momentum.Rsi[float64]

	// Sma represents the configuration parameters for calculating the Simple Moving Average (SMA).
	Sma *trend.Sma[float64]

	// DownDays is the number of down days for RSI.
	DownDays int

	// BuySignalAt defines the RSI level at which a Buy signal is confirmed.
	BuySignalAt float64

	// BuyAt defines the RSI level at which a Buy action is generated.
	BuyAt float64

	// SellAt defines the RSI level at which a Sell action is generated.
	SellAt float64
}

// NewTripleRsiStrategy function initializes a new Triple RSI strategy instance with the default parameters.
func NewTripleRsiStrategy() *TripleRsiStrategy {
	return NewTripleRsiStrategyWith(
		DefaultTripleRsiStrategyPeriod,
		DefaultTripleRsiStrategyMovingAveragePeriod,
		DefaultTripleRsiStrategyDownDays,
		DefaultTripleRsiStrategyBuySignalAt,
		DefaultTripleRsiStrategyBuyAt,
		DefaultTripleRsiStrategySellAt,
	)
}

// NewTripleRsiStrategyWith function initializes a new RSI strategy instance with the given parameters.
func NewTripleRsiStrategyWith(period, smaPeriod, downDays int, buySignalAt, buyAt, sellAt float64) *TripleRsiStrategy {
	return &TripleRsiStrategy{
		Rsi:         momentum.NewRsiWithPeriod[float64](period),
		Sma:         trend.NewSmaWithPeriod[float64](smaPeriod),
		DownDays:    downDays,
		BuySignalAt: buySignalAt,
		BuyAt:       buyAt,
		SellAt:      sellAt,
	}
}

// Name returns the name of the strategy.
func (t *TripleRsiStrategy) Name() string {
	return fmt.Sprintf("Triple RSI Strategy (%d,%d,%d,%.0f,%.0f,%.0f)", t.Rsi.Rma.Period, t.Sma.Period, t.DownDays,
		t.BuySignalAt, t.BuyAt, t.SellAt)
}

// IdlePeriod is the initial period that the Triple RSI strategy won't yield any results.
func (t *TripleRsiStrategy) IdlePeriod() int {
	return t.Sma.IdlePeriod()
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (t *TripleRsiStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots),
		3,
	)

	rsis := t.Rsi.Compute(closingsSplice[0])
	smas := t.Sma.Compute(closingsSplice[1])
	memory := helper.NewRing[float64](t.DownDays)

	// Skip RSI results until SMA is ready.
	rsis = helper.Skip(rsis, t.Sma.IdlePeriod()-t.Rsi.IdlePeriod())

	// Skip closing values until SMA is ready.
	closingsSplice[2] = helper.Skip(closingsSplice[2], t.Sma.IdlePeriod())

	actions := helper.Operate3(rsis, smas, closingsSplice[2], func(rsi, sma, closing float64) strategy.Action {
		memory.Put(rsi)

		if !memory.IsFull() {
			return strategy.Hold
		}

		// Recommend Sell:
		// - Sell at the close when the 5-period RSI crosses above 50.
		if rsi > t.SellAt {
			return strategy.Sell
		}

		// Recommend Buy:
		// - The 5-period RSI is below 30.
		if rsi >= t.BuyAt {
			return strategy.Hold
		}

		// - The 5-period RSI reading is down for the 3rd period in a row.
		for i := 1; i < t.DownDays; i++ {
			if memory.At(i-1) > memory.At(i) {
				return strategy.Hold
			}
		}

		// - The 5-period RSI reading was below 60 three trading periods ago.
		if memory.At(0) >= t.BuySignalAt {
			return strategy.Hold
		}

		// - The close is higher than the 200-period moving average.
		if closing <= sma {
			return strategy.Hold
		}

		return strategy.Buy
	})

	// Shift actions until strategy is ready.
	actions = helper.Shift(actions, t.Sma.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (t *TripleRsiStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> Compute     -> actions -> annotations
	// snapshots[2] -> closings[0] -> close
	//              -> closings[1] -> Rsi.Compute -> rsi
	//              -> closings[2] -> Sma.Compute -> sma
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[2]), 3)

	rsis := t.Rsi.Compute(closings[1])
	smas := t.Sma.Compute(closings[2])

	actions, outcomes := strategy.ComputeWithOutcome(t, snapshots[1])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	dates = helper.Skip(dates, t.IdlePeriod())
	closings[0] = helper.Skip(closings[0], t.IdlePeriod())
	rsis = helper.Skip(rsis, t.IdlePeriod()-t.Rsi.IdlePeriod())
	annotations = helper.Skip(annotations, t.IdlePeriod())
	outcomes = helper.Skip(outcomes, t.IdlePeriod())

	report := helper.NewReport(t.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn(fmt.Sprintf("RSI(%d)", t.Rsi.Rma.Period), rsis), 1)
	report.AddColumn(helper.NewNumericReportColumn(fmt.Sprintf("SMA(%d)", t.Sma.Period), smas))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)
	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
