// Package backtest contains the backtest functions.
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
package backtest

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultBacktestWorkers is the default number of backtest workers.
	DefaultBacktestWorkers = 1

	// DefaultLastDays is the default number of days backtest should go back.
	DefaultLastDays = 365
)

// Backtest function rigorously evaluates the potential performance of the
// specified strategies applied to a defined set of assets. It generates
// comprehensive visual representations for each strategy-asset pairing.
type Backtest struct {
	// repository is the repository to retrieve the assets from.
	repository asset.Repository

	// report is the report writer for the backtest.
	report Report

	// Names is the names of the assets to backtest.
	Names []string

	// Strategies is the list of strategies to apply.
	Strategies []strategy.Strategy

	// Workers is the number of concurrent workers.
	Workers int

	// LastDays is the number of days backtest should go back.
	LastDays int
}

// NewBacktest function initializes a new backtest instance.
func NewBacktest(repository asset.Repository, report Report) *Backtest {
	return &Backtest{
		repository: repository,
		report:     report,
		Names:      []string{},
		Strategies: []strategy.Strategy{},
		Workers:    DefaultBacktestWorkers,
		LastDays:   DefaultLastDays,
	}
}

// Run executes a comprehensive performance evaluation of the designated strategies,
// applied to a specified collection of assets. In the absence of explicitly defined
// assets, encompasses all assets within the repository. Likewise, in the absence of
// explicitly defined strategies, encompasses all the registered strategies.
func (b *Backtest) Run() error {
	// When asset names are absent, considers all assets within the provided repository for evaluation.
	if len(b.Names) == 0 {
		assets, err := b.repository.Assets()
		if err != nil {
			return err
		}

		b.Names = assets
	}

	// When strategies are absent, considers all strategies.
	if len(b.Strategies) == 0 {
		b.Strategies = []strategy.Strategy{
			strategy.NewBuyAndHoldStrategy(),
		}
	}

	// Begin report.
	err := b.report.Begin(b.Names, b.Strategies)
	if err != nil {
		return fmt.Errorf("unable to begin report: %w", err)
	}

	// Run the backtest workers.
	names := helper.SliceToChan(b.Names)
	wg := &sync.WaitGroup{}

	for i := 0; i < b.Workers; i++ {
		wg.Add(1)
		go b.worker(names, wg)
	}

	// Wait for all workers to finish.
	wg.Wait()

	// End report.
	err = b.report.End()
	if err != nil {
		return fmt.Errorf("unable to end report: %w", err)
	}

	return nil
}

// worker is a backtesting worker that concurrently executes backtests for individual
// assets. It receives asset names from the provided channel, and performs backtests
// using the given strategies.
func (b *Backtest) worker(names <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	since := time.Now().AddDate(0, 0, -b.LastDays)

	for name := range names {
		log.Printf("Backtesting %s...", name)
		snapshots, err := b.repository.GetSince(name, since)
		if err != nil {
			log.Printf("Unable to retrieve the snapshots for %s: %v", name, err)
			continue
		}

		// We don't expect the snapshots to be a stream during backtesting.
		snapshotsSlice := helper.ChanToSlice(snapshots)

		// Backtesting asset has begun.
		err = b.report.AssetBegin(name, b.Strategies)
		if err != nil {
			log.Printf("Unable to asset begin for %s: %v", name, err)
			continue
		}

		// Backtest strategies on the given asset.
		for _, currentStrategy := range b.Strategies {
			snapshotsSplice := helper.Duplicate(helper.SliceToChan(snapshotsSlice), 2)

			actions, outcomes := strategy.ComputeWithOutcome(currentStrategy, snapshotsSplice[0])
			err = b.report.Write(name, currentStrategy, snapshotsSplice[1], actions, outcomes)
			if err != nil {
				log.Printf("Unable to write report for %s: %v", name, err)
			}
		}

		// Backtesting asset had ended
		err = b.report.AssetEnd(name)
		if err != nil {
			log.Printf("Unable to asset end for %s: %v", name, err)
		}
	}
}
