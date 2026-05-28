// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// AndStrategy combines multiple strategies and emits actionable recommendations when **all** strategies in
// the group **reach the same actionable conclusion**. This can be a conservative approach, potentially
// delaying recommendations until full consensus is reached.
type AndStrategy struct {
	// Strategies are the group of strategies that will be consulted to make an actionable recommendation.
	Strategies []Strategy

	// name is the name of this group of strategies.
	name string
}

// NewAndStrategy function initializes an empty and strategies group with the given name.
func NewAndStrategy(name string, strategies ...Strategy) *AndStrategy {
	return &AndStrategy{
		Strategies: strategies,
		name:       name,
	}
}

// Name returns the name of the strategy.
func (a *AndStrategy) Name() string {
	return a.name
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (a *AndStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan Action {
	result := make(chan Action)

	sources := ActionSources(a.Strategies, snapshots)

	go func() {
		defer close(result)

		and := len(a.Strategies)

		for {
			buy, _, sell, ok := CountActions(sources)
			if !ok {
				break
			}

			if sell == and {
				result <- Sell
			} else if buy == and {
				result <- Buy
			} else {
				result <- Hold
			}
		}
	}()

	return result
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (a *AndStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
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

// AllAndStrategies performs a cartesian product operation on the given strategies, resulting in a collection
// containing all and strategies formed by combining two strategies together.
func AllAndStrategies(strategies []Strategy) []Strategy {
	andStrategies := make([]Strategy, 0, len(strategies)*len(strategies))

	for _, first := range strategies {
		for _, second := range strategies {
			if first != second {
				andStrategy := NewAndStrategy(fmt.Sprintf("%s and %s", first.Name(), second.Name()))
				andStrategy.Strategies = []Strategy{first, second}
				andStrategies = append(andStrategies, andStrategy)
			}
		}
	}

	return andStrategies
}
