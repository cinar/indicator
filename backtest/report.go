// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest

import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/strategy"
)

// Report is the backtest report interface.
type Report interface {
	// Begin is called when the backtest begins.
	Begin(assetNames []string, strategies []strategy.Strategy) error

	// AssetBegin is called when backtesting for the given asset begins.
	AssetBegin(name string, strategies []strategy.Strategy) error

	// Write writes the given strategy actions and outomes to the report.
	Write(assetName string, currentStrategy strategy.Strategy, snapshots <-chan *asset.Snapshot, actions <-chan strategy.Action, outcomes <-chan float64) error

	// AssetEnd is called when backtesting for the given asset ends.
	AssetEnd(name string) error

	// End is called when the backtest ends.
	End() error
}
