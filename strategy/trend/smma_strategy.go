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

const (
	// DefaultSmmaStrategyShortPeriod is the default short-term SMMA period of 20.
	DefaultSmmaStrategyShortPeriod = 20

	// DefaultSmmaStrategyLongPeriod is the default short-term SMMA period of 50.
	DefaultSmmaStrategyLongPeriod = 50
)

// SmmaStrategy represents the configuration parameters for calculating the
// Smooted Moving Averge (SMMA) strategy. A short-term SMMA crossing above
// the long-term SMMA suggests a bullish trend, while crossing below the
// long-term SMMA indicates a bearish trend.
type SmmaStrategy struct {
	// ShortSmma represents the configuration parameters for calculating the
	// short-term Smooted Moving Averge (SMMA).
	ShortSmma *trend.Smma[float64]

	// LongSmma represents the configuration parameters for calculating the
	// long-term Smooted Moving Averge (SMMA).
	LongSmma *trend.Smma[float64]
}

// NewSmmaStrategy function initializes a new SMMA strategy instance.
func NewSmmaStrategy() *SmmaStrategy {
	return NewSmmaStrategyWith(
		DefaultSmmaStrategyShortPeriod,
		DefaultSmmaStrategyLongPeriod,
	)
}

// NewSmmaStrategyWith function initializes a new SMMA strategy instance with the given parameters.
func NewSmmaStrategyWith(shortPeriod, longPeriod int) *SmmaStrategy {
	return &SmmaStrategy{
		ShortSmma: trend.NewSmmaWithPeriod[float64](shortPeriod),
		LongSmma:  trend.NewSmmaWithPeriod[float64](longPeriod),
	}
}

// Name returns the name of the strategy.
func (s *SmmaStrategy) Name() string {
	return fmt.Sprintf("SMMA Strategy (%d,%d)",
		s.ShortSmma.Period,
		s.LongSmma.Period,
	)
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *SmmaStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closingsSplice := helper.Duplicate(asset.SnapshotsAsClosings(snapshots), 2)

	shortSmmas := s.ShortSmma.Compute(closingsSplice[0])
	longSmmas := s.LongSmma.Compute(closingsSplice[1])

	commonPeriod := helper.CommonPeriod(s.ShortSmma.Period, s.LongSmma.Period)
	shortSmmas = helper.SyncPeriod(commonPeriod, s.ShortSmma.Period, shortSmmas)
	longSmmas = helper.SyncPeriod(commonPeriod, s.LongSmma.Period, longSmmas)

	actions := helper.Operate(shortSmmas, longSmmas, func(shortSmma, longSmma float64) strategy.Action {
		// A short-perios SMMA value crossing above long-period SMMA suggests a bullish trend.
		if shortSmma > longSmma {
			return strategy.Buy
		}

		// A short-period SMMA value crossing below long-period SMMA suggests a bearish trend.
		if longSmma > shortSmma {
			return strategy.Sell
		}

		return strategy.Hold
	})

	// SMMA strategy starts only after a full period.
	actions = helper.Shift(actions, commonPeriod, strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (s *SmmaStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> short-period SMMA
	//                 closings[2] -> long-period SMMA
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 3)

	shortSmmas := s.ShortSmma.Compute(closings[1])
	longSmmas := s.LongSmma.Compute(closings[2])

	actions, outcomes := strategy.ComputeWithOutcome(s, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	commonPeriod := helper.CommonPeriod(s.ShortSmma.Period, s.LongSmma.Period)
	dates = helper.SyncPeriod(commonPeriod, 0, dates)
	closings[0] = helper.Skip(closings[0], commonPeriod)
	shortSmmas = helper.SyncPeriod(commonPeriod, s.ShortSmma.Period, shortSmmas)
	longSmmas = helper.SyncPeriod(commonPeriod, s.LongSmma.Period, longSmmas)
	annotations = helper.Skip(annotations, commonPeriod)
	outcomes = helper.Skip(outcomes, commonPeriod)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[0]))
	report.AddColumn(helper.NewNumericReportColumn("MACD", shortSmmas), 1)
	report.AddColumn(helper.NewNumericReportColumn("Signal", longSmmas), 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 2)

	return report
}
