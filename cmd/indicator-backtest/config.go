// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/strategy"
)

// SourceConfig names a registered repository or report builder along with the configuration
// string it is built with, mirroring the -repository-name/-repository-config and
// -report-name/-report-config flag pairs.
type SourceConfig struct {
	// Name is the registered repository or report builder name.
	Name string `json:"name"`

	// Config is the configuration string passed to the builder.
	Config string `json:"config"`
}

// StrategyConfig names a registered strategy along with an optional JSON fragment that is
// overlaid onto the strategy's default instance. See NewStrategyFromConfig for how the
// overlay is applied.
type StrategyConfig struct {
	// Name is the registered strategy name (see StrategyNames).
	Name string `json:"name"`

	// Config is an optional JSON object overlaid onto the strategy's exported fields, letting
	// any indicator parameter the strategy exposes (e.g. a period, or a buy/sell threshold) be
	// tuned without changing the strategy's registered defaults.
	Config json.RawMessage `json:"config,omitempty"`
}

// Config is the top level structure for a JSON indicator-backtest configuration file. It
// mirrors the command line flags, letting a full backtest run, including every strategy's
// parameters, be described in a single file instead of a long flag list.
type Config struct {
	// Repository configures the asset repository the backtest reads snapshots from.
	Repository SourceConfig `json:"repository"`

	// Report configures the report writer the backtest results are written to.
	Report SourceConfig `json:"report"`

	// Symbols is the list of asset names to backtest. When empty, all assets in the
	// repository are used.
	Symbols []string `json:"symbols,omitempty"`

	// Strategies is the list of strategies to run the backtest with.
	Strategies []StrategyConfig `json:"strategies"`

	// Workers is the number of concurrent workers. Defaults to backtest.DefaultBacktestWorkers
	// when zero.
	Workers int `json:"workers,omitempty"`

	// LastDays is the number of days the backtest should go back. Defaults to
	// backtest.DefaultLastDays when zero.
	LastDays int `json:"lastDays,omitempty"`

	// AddSplits adds the cartesian product of split strategies built from Strategies.
	AddSplits bool `json:"addSplits,omitempty"`

	// AddAnds adds the cartesian product of and strategies built from Strategies.
	AddAnds bool `json:"addAnds,omitempty"`
}

// LoadConfig reads and parses the JSON indicator-backtest configuration file at the given path.
func LoadConfig(path string) (*Config, error) {
	// path is an operator-supplied CLI argument (-config), the same trust level as every
	// other -*-config flag this command already accepts, not untrusted network input.
	data, err := os.ReadFile(path) // #nosec G304 -- path is the local file named by the -config flag.
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	config := &Config{
		Repository: SourceConfig{Name: "filesystem"},
		Report:     SourceConfig{Name: "html", Config: "."},
		Workers:    backtest.DefaultBacktestWorkers,
		LastDays:   backtest.DefaultLastDays,
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %w", err)
	}

	if len(config.Strategies) == 0 {
		return nil, fmt.Errorf("config file %s does not name any strategies", path)
	}

	return config, nil
}

// NewBacktestFromConfig builds a ready-to-run backtest.Backtest from the given Config.
func NewBacktestFromConfig(config *Config, logger *slog.Logger) (*backtest.Backtest, error) {
	source, err := asset.NewRepository(config.Repository.Name, config.Repository.Config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize source: %w", err)
	}

	report, err := backtest.NewReport(config.Report.Name, config.Report.Config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize report: %w", err)
	}

	backtester := backtest.NewBacktest(source, report)
	backtester.Workers = config.Workers
	backtester.LastDays = config.LastDays
	backtester.Logger = logger
	backtester.Names = append(backtester.Names, config.Symbols...)

	for _, sc := range config.Strategies {
		s, err := NewStrategyFromConfig(sc.Name, sc.Config)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize strategy %q: %w", sc.Name, err)
		}

		backtester.Strategies = append(backtester.Strategies, s)
	}

	if config.AddSplits {
		backtester.Strategies = append(backtester.Strategies, strategy.AllSplitStrategies(backtester.Strategies)...)
	}

	if config.AddAnds {
		backtester.Strategies = append(backtester.Strategies, strategy.AllAndStrategies(backtester.Strategies)...)
	}

	return backtester, nil
}
