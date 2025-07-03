package main

import (
	"errors"
	"fmt"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

// StrategyRequest represents the request structure for strategy processing
type StrategyRequest struct {
	Strategy StrategyType `json:"strategy"`
	Data     OhlcvData    `json:"data"`
}

// OhlcvData represents the OHLCV data structure for the strategy
type OhlcvData struct {
	Date    []int64   `json:"date"`
	Opening []float64 `json:"opening"`
	Closing []float64 `json:"closing"`
	High    []float64 `json:"high"`
	Low     []float64 `json:"low"`
	Volume  []float64 `json:"volume"`
}

// Response represents the JSON response structure
type Response struct {
	Actions []strategy.Action `json:"actions"`
}

// runBacktest processes the OHLCV data using the specified strategy and returns the actions
func runBacktest(strategyType StrategyType, data OhlcvData) ([]*backtest.DataStrategyResult, error) {
	dateArray := toTimeArray(data.Date)
	snapshots := make(chan *asset.Snapshot)

	go func() {
		defer close(snapshots)

		for i := range dateArray {
			snapshots <- &asset.Snapshot{
				Date:   dateArray[i],
				Open:   data.Opening[i],
				Close:  data.Closing[i],
				High:   data.High[i],
				Low:    data.Low[i],
				Volume: data.Volume[i],
			}
		}
	}()

	repository := asset.NewInMemoryRepository()
	err := repository.Append("in_memory_asset", snapshots)
	if err != nil {
		return nil, fmt.Errorf("failed to append snapshots: %w", err)
	}

	// Create the strategy using the factory function
	strat, err := CreateStrategy(strategyType)
	if err != nil {
		return nil, fmt.Errorf("failed to create strategy: %w", err)
	}

	assets := []string{"in_memory_asset"}
	strategies := []strategy.Strategy{strat}

	report := backtest.NewDataReport()

	err = report.Begin(assets, strategies)
	if err != nil {
		return nil, err
	}

	err = report.AssetBegin(assets[0], strategies)
	if err != nil {
		return nil, err
	}

	snapshotsChan, err := repository.Get(assets[0])
	if err != nil {
		return nil, err
	}

	// Duplicate the snapshots channel 3 ways to be used for:
	// 1. Report generation (snapshotsSplice[0])
	// 2. Strategy computation (snapshotsSplice[1])
	// 3. Outcome calculation (snapshotsSplice[2])
	snapshotsSplice := helper.Duplicate(snapshotsChan, 3)

	// Compute strategy actions and duplicate the results 2 ways to be used for:
	// 1. Report generation (actionsSplice[0])
	// 2. Outcome calculation (actionsSplice[1])
	actionsSplice := helper.Duplicate(
		strategies[0].Compute(snapshotsSplice[1]),
		2,
	)

	// Calculate outcomes using the third copy of snapshots (converted to closing prices)
	// and the second copy of the computed actions
	outcomes := strategy.Outcome(
		asset.SnapshotsAsClosings(snapshotsSplice[2]),
		actionsSplice[1],
	)

	err = report.Write(assets[0], strategies[0], snapshotsSplice[0], actionsSplice[0], outcomes)
	if err != nil {
		return nil, err
	}

	err = report.AssetEnd(assets[0])
	if err != nil {
		return nil, err
	}

	err = report.End()
	if err != nil {
		return nil, err
	}

	results, ok := report.Results[assets[0]]
	if !ok {
		return nil, errors.New("asset result not found")
	}

	if len(results) != len(strategies) {
		return nil, fmt.Errorf("results count and strategies count are not the same, %d %d", len(results), len(strategies))
	}

	return results, nil
}
