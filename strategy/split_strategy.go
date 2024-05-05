// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// SplitStrategy leverages two separate strategies. It utilizes the first strategy to identify potential Buy
// opportunities, and the second strategy to identify potential Sell opportunities. When there is a
// conflicting recommendation, returns Hold.
type SplitStrategy struct {
	// BuyStrategy is used to identify potential Buy opportunities.
	BuyStrategy Strategy

	// SellStrategy is used to identify potential Sell opportunities.
	SellStrategy Strategy
}

// NewSplitStrategy function initializes a new split strategy with the given parameters.
func NewSplitStrategy(buyStrategy, sellStrategy Strategy) *SplitStrategy {
	return &SplitStrategy{
		BuyStrategy:  buyStrategy,
		SellStrategy: sellStrategy,
	}
}

// Name returns the name of the strategy.
func (s *SplitStrategy) Name() string {
	return fmt.Sprintf("SplitStrategy(%s, %s)", s.BuyStrategy.Name(), s.SellStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (s *SplitStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan Action {
	result := make(chan Action)

	snapshotsSplice := helper.Duplicate(snapshots, 2)

	buyActions := s.BuyStrategy.Compute(snapshotsSplice[0])
	sellActions := s.SellStrategy.Compute(snapshotsSplice[1])

	go func() {
		defer close(result)

		for {
			buyAction, ok := <-buyActions
			if !ok {
				break
			}

			sellAction, ok := <-sellActions
			if !ok {
				break
			}

			if (buyAction == Buy) && (sellAction != Sell) {
				result <- Buy
			} else if (sellAction == Sell) && (buyAction != Buy) {
				result <- Sell
			} else {
				result <- Hold
			}
		}
	}()

	return NormalizeActions(result)
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (s *SplitStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := ComputeWithOutcome(s, snapshots[2])
	annotations := ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(s.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// AllSplitStrategies performs a cartesian product operation on the given strategies, resulting in a collection
// containing all split strategies formed by combining individual buy and sell strategies.
func AllSplitStrategies(strategies []Strategy) []Strategy {
	splitStrategies := make([]Strategy, 0, len(strategies)*len(strategies))

	for _, buyStrategy := range strategies {
		for _, sellStrategy := range strategies {
			if buyStrategy != sellStrategy {
				splitStrategies = append(splitStrategies, NewSplitStrategy(buyStrategy, sellStrategy))
			}
		}
	}

	return splitStrategies
}
