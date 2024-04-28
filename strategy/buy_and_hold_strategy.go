// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// BuyAndHoldStrategy defines an investment approach of acquiring and
// indefinitely retaining an asset. This strategy primarily serves as
// a benchmark for evaluating the performance of alternative
// strategies against a baseline of passive asset ownership.
type BuyAndHoldStrategy struct {
	Strategy
}

// NewBuyAndHoldStrategy function initializes a new buy and hold strategy instance.
func NewBuyAndHoldStrategy() *BuyAndHoldStrategy {
	return &BuyAndHoldStrategy{}
}

// Name returns the name of the strategy.
func (*BuyAndHoldStrategy) Name() string {
	return "Buy and Hold Strategy"
}

// Compute processes the provided asset snapshots and generates a
// stream of actionable recommendations.
func (*BuyAndHoldStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan Action {
	closings := asset.SnapshotsAsClosings(snapshots)
	actions := make(chan Action, cap(snapshots))

	go func() {
		defer close(actions)

		_, ok := <-closings
		if !ok {
			return
		}

		actions <- Buy

		for range closings {
			actions <- Hold
		}
	}()

	return actions
}

// Report processes the provided asset snapshots and generates a
// report annotated with the recommended actions.
func (b *BuyAndHoldStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := ComputeWithOutcome(b, snapshots[2])
	annotations := ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(b.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
