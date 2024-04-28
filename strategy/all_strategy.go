// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// AllStrategy combines multiple strategies and emits actionable recommendations when **all** strategies in
// the group **reach the same actionable conclusion**. This can be a conservative approach, potentially
// delaying recommendations until full consensus is reached.
type AllStrategy struct {
	Strategy

	// Strategies are the group of strategies that will be consulted to make an actionable recommendation.
	Strategies []Strategy

	// name is the name of this group of strategies.
	name string
}

// NewAllStrategy function initializes an empty all strategies group with the given name.
func NewAllStrategy(name string) *AllStrategy {
	return &AllStrategy{
		Strategies: []Strategy{},
		name:       name,
	}
}

// Name returns the name of the strategy.
func (a *AllStrategy) Name() string {
	return a.name
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (a *AllStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan Action {
	result := make(chan Action)

	sources := ActionSources(a.Strategies, snapshots)

	go func() {
		defer close(result)

		all := len(a.Strategies)

		for {
			buy, _, sell, ok := CountActions(sources)
			if !ok {
				break
			}

			if sell == all {
				result <- Sell
			} else if buy == all {
				result <- Buy
			} else {
				result <- Hold
			}
		}
	}()

	return NormalizeActions(result)
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (a *AllStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := ComputeWithOutcome(a, snapshots[2])
	annotations := ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(a.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
