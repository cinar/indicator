// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// DataStrategyResult is the strategy result.
type DataStrategyResult struct {
	// Asset is the asset name.
	Asset string

	// Strategy is the strategy instnace.
	Strategy strategy.Strategy

	// Outcome is the strategy outcome.
	Outcome float64

	// Action is the final action recommended by the strategy.
	Action strategy.Action

	// Transactions are the action recommendations.
	Transactions []strategy.Action
}

// DataReport is the bactest data report enablign programmatic access to the backtest results.
type DataReport struct {
	// Results are the backtest results for the assets.
	Results map[string][]*DataStrategyResult
}

// NewDataReport initializes a new data report instance.
func NewDataReport() *DataReport {
	return &DataReport{
		Results: make(map[string][]*DataStrategyResult),
	}
}

// Begin is called when the backtest begins.
func (*DataReport) Begin(_ []string, _ []strategy.Strategy) error {
	return nil
}

// AssetBegin is called when backtesting for the given asset begins.
func (d *DataReport) AssetBegin(name string, strategies []strategy.Strategy) error {
	d.Results[name] = make([]*DataStrategyResult, 0, len(strategies))
	return nil
}

// Write writes the given strategy actions and outomes to the report.
func (d *DataReport) Write(assetName string, currentStrategy strategy.Strategy, snapshots <-chan *asset.Snapshot, actions <-chan strategy.Action, outcomes <-chan float64) error {
	go helper.Drain(snapshots)

	actionsSplice := helper.Duplicate(actions, 2)

	lastOutcome := helper.Last(outcomes, 1)
	lastAction := helper.Last(actionsSplice[0], 1)
	transactions := helper.ChanToSlice(actionsSplice[1])

	result := &DataStrategyResult{
		Asset:        assetName,
		Strategy:     currentStrategy,
		Outcome:      <-lastOutcome,
		Action:       <-lastAction,
		Transactions: transactions,
	}

	d.Results[assetName] = append(d.Results[assetName], result)

	return nil
}

// AssetEnd is called when backtesting for the given asset ends.
func (*DataReport) AssetEnd(_ string) error {
	return nil
}

// End is called when the backtest ends.
func (*DataReport) End() error {
	return nil
}
