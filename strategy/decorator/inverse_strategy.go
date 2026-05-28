// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package decorator

import (
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// InverseStrategy reverses the advice of another strategy. For example, if the original strategy suggests buying an
// asset, InverseStrategy would recommend selling it.
type InverseStrategy struct {
	// InnerStrategy is the inner strategy.
	InnerStrategy strategy.Strategy
}

// NewInverseStrategy function initializes a new inverse strategy instance.
func NewInverseStrategy(innerStrategy strategy.Strategy) *InverseStrategy {
	return &InverseStrategy{
		InnerStrategy: innerStrategy,
	}
}

// Name returns the name of the strategy.
func (i *InverseStrategy) Name() string {
	return fmt.Sprintf("Inverse Strategy (%s)", i.InnerStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (i *InverseStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	return helper.Map(i.InnerStrategy.Compute(snapshots), func(action strategy.Action) strategy.Action {
		switch action {
		case strategy.Buy:
			return strategy.Sell

		case strategy.Sell:
			return strategy.Buy

		default:
			return strategy.Hold
		}
	})
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (i *InverseStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(i, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(i.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
