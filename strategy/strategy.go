// Package strategy contains the strategy functions.
//
// This package belongs to the Indicator project. Indicator is
// a Golang module that supplies a variety of technical
// indicators, strategies, and a backtesting framework
// for analysis.
//
// # License
//
//	Copyright (c) 2021-2024 Onur Cinar.
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
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// Strategy defines a shared interface for trading strategies.
type Strategy interface {
	// Name returns the name of the strategy.
	Name() string

	// Compute processes the provided asset snapshots and generates a
	// stream of actionable recommendations.
	Compute(snapshots <-chan *asset.Snapshot) <-chan Action

	// Report processes the provided asset snapshots and generates a
	// report annotated with the recommended actions.
	Report(snapshots <-chan *asset.Snapshot) *helper.Report
}

// ComputeWithOutcome uses the given strategy to processes the provided asset snapshots and
// generates a stream of actionable recommendations and outcomes.
func ComputeWithOutcome(s Strategy, c <-chan *asset.Snapshot) (<-chan Action, <-chan float64) {
	snapshots := helper.Duplicate(c, 2)

	actions := helper.Duplicate(s.Compute(snapshots[0]), 2)

	openings := helper.Skip(asset.SnapshotsAsOpenings(snapshots[1]), 1)

	outcomes := helper.Echo(
		Outcome(openings, actions[0]),
		1,
		1,
	)

	return actions[1], outcomes
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
