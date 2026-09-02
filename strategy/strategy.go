// Package strategy contains the strategy functions.
//
// This package belongs to the Indicator project. Indicator is
// a Golang module that supplies a variety of technical
// indicators, strategies, and a backtesting framework
// for analysis.
//
// # License
//
//	Copyright (c) 2021-2026 The Indicator Authors.
//	The source code is provided under GNU AGPLv3 License.
//	https://github.com/cinar/indicator
//
// # Disclaimer
//
// The information provided on this project is strictly for
// informational purposes and is not to be construed as
// advice or solicitation to buy or sell any security.
package strategy

import (
	"context"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// Strategy defines a shared interface for trading strategies.
type Strategy interface {
	// Name returns the name of the example strategy.
	Name() string

	// Compute processes the provided asset snapshots and generates a
	// stream of actionable recommendations.
	Compute(snapshots <-chan *asset.Snapshot) <-chan Action

	// Report processes the provided asset snapshots and generates a
	// report annotated with the recommended actions.
	Report(snapshots <-chan *asset.Snapshot) *helper.Report
}

// StrategyWithContext defines a shared interface for trading strategies
// supporting context-aware computations.
type StrategyWithContext interface {
	Strategy
	ComputeWithContext(ctx context.Context, snapshots <-chan *asset.Snapshot) <-chan Action
}

// ComputeStrategyWithContext processes snapshots with a strategy using context.
func ComputeStrategyWithContext(ctx context.Context, s Strategy, c <-chan *asset.Snapshot) <-chan Action {
	if sc, ok := s.(StrategyWithContext); ok {
		return sc.ComputeWithContext(ctx, c)
	}
	return s.Compute(c)
}

// ComputeWithOutcomeWithContext uses the given strategy to processes the provided asset snapshots and
// generates a stream of actionable recommendations and outcomes, supporting context cancellation.
func ComputeWithOutcomeWithContext(ctx context.Context, s Strategy, c <-chan *asset.Snapshot) (<-chan Action, <-chan float64) {
	snapshots := helper.DuplicateWithContext(ctx, c, 2)

	actions := helper.DuplicateWithContext(ctx, ComputeStrategyWithContext(ctx, s, snapshots[0]), 2)
	closings := asset.SnapshotsAsClosingsWithContext(ctx, snapshots[1])

	outcomes := OutcomeWithContext(ctx, closings, actions[1])

	return actions[0], outcomes
}

// ComputeWithOutcome uses the given strategy to processes the provided asset snapshots and
// generates a stream of actionable recommendations and outcomes.
//
// Deprecated: Use ComputeWithOutcomeWithContext instead.
func ComputeWithOutcome(s Strategy, c <-chan *asset.Snapshot) (<-chan Action, <-chan float64) {
	return ComputeWithOutcomeWithContext(context.Background(), s, c)
}

// ComputeWithOutcomeAndTimingWithContext uses the given strategy to process the provided asset snapshots and
// generates a stream of actionable recommendations and outcomes, using the given ExecutionTiming to decide
// which price a simulated trade executes at, supporting context cancellation.
//
// With AtClose, this behaves identically to ComputeWithOutcomeWithContext: each action is paired with the
// closing price of the same bar it was computed from.
//
// With NextOpen or NextClose, each action is instead paired with the opening or closing price of the
// following bar. As a result, the returned outcomes channel yields one fewer value than the returned
// actions channel, since there is no following bar for the last action. Callers that need the actions and
// outcomes channels aligned position-for-position (for example, to build a report column) must account for
// this offset themselves, the same way many strategies already skip-align channels of differing lengths.
// Callers only interested in the final/aggregate outcome (for example, via helper.Last(outcomes, 1)) are
// unaffected either way.
func ComputeWithOutcomeAndTimingWithContext(ctx context.Context, s Strategy, c <-chan *asset.Snapshot, timing ExecutionTiming) (<-chan Action, <-chan float64) {
	snapshots := helper.DuplicateWithContext(ctx, c, 2)

	actions := helper.DuplicateWithContext(ctx, ComputeStrategyWithContext(ctx, s, snapshots[0]), 2)

	var prices <-chan float64

	switch timing {
	case NextOpen:
		prices = helper.SkipWithContext(ctx, asset.SnapshotsAsOpeningsWithContext(ctx, snapshots[1]), 1)

	case NextClose:
		prices = helper.SkipWithContext(ctx, asset.SnapshotsAsClosingsWithContext(ctx, snapshots[1]), 1)

	default:
		prices = asset.SnapshotsAsClosingsWithContext(ctx, snapshots[1])
	}

	outcomes := OutcomeWithContext(ctx, prices, actions[1])

	return actions[0], outcomes
}

// ComputeWithOutcomeAndTiming uses the given strategy to process the provided asset snapshots and
// generates a stream of actionable recommendations and outcomes, using the given ExecutionTiming to decide
// which price a simulated trade executes at.
//
// See ComputeWithOutcomeAndTimingWithContext for details, including the note on the outcomes channel
// being shorter than the actions channel when timing is NextOpen or NextClose.
func ComputeWithOutcomeAndTiming(s Strategy, c <-chan *asset.Snapshot, timing ExecutionTiming) (<-chan Action, <-chan float64) {
	return ComputeWithOutcomeAndTimingWithContext(context.Background(), s, c, timing)
}

// AllStrategies returns a slice containing references to all available base strategies.
func AllStrategies() []Strategy {
	return []Strategy{
		NewBuyAndHoldStrategy(),
	}
}

// ActionSources creates a slice of action channels, one for each strategy, where each channel emits actions
// computed by its corresponding strategy based on snapshots from the provided snapshot channel.
func ActionSources(strategies []Strategy, snapshots <-chan *asset.Snapshot) []<-chan Action {
	snapshotsSplice := helper.Duplicate(snapshots, len(strategies))
	sources := make([]<-chan Action, len(strategies))

	for i, strategy := range strategies {
		sources[i] = DenormalizeActions(
			strategy.Compute(snapshotsSplice[i]),
		)
	}

	return sources
}
