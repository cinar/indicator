// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	// Go embed report template.
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sync"
	"text/template"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultBacktestWorkers is the default number of backtest workers.
	DefaultBacktestWorkers = 1

	// DefaultLastDays is the default number of days backtest should go back.
	DefaultLastDays = 30
)

//go:embed "backtest_report.tmpl"
var backtestReportTmpl string

//go:embed "backtest_asset_report.tmpl"
var backtestAssetReportTmpl string

// Backtest function rigorously evaluates the potential performance of the
// specified strategies applied to a defined set of assets. It generates
// comprehensive visual representations for each strategy-asset pairing.
type Backtest struct {
	// repository is the repository to retrieve the assets from.
	repository asset.Repository

	// outputDir is the output directory for the generated reports.
	outputDir string

	// Names is the names of the assets to backtest.
	Names []string

	// Strategies is the list of strategies to apply.
	Strategies []Strategy

	// Workers is the number of concurrent workers.
	Workers int

	// LastDays is the number of days backtest should go back.
	LastDays int
}

// backtestResult encapsulates the outcome of running a strategy.
type backtestResult struct {
	// AssetName is the name of the asset.
	AssetName string

	// StrategyName is the name of the strategy.
	StrategyName string

	// Action is the last recommended action by the strategy.
	Action Action

	// Since indicates how long the current action recommendation has been in effect.
	Since int

	// Outcome is the effectiveness of applying the recommended actions.
	Outcome float64

	// Transactions is the number of transactions made by the strategy.
	Transactions int
}

// NewBacktest function initializes a new backtest instance.
func NewBacktest(repository asset.Repository, outputDir string) *Backtest {
	return &Backtest{
		repository: repository,
		outputDir:  outputDir,
		Names:      []string{},
		Strategies: []Strategy{},
		Workers:    DefaultBacktestWorkers,
		LastDays:   DefaultLastDays,
	}
}

// Run executes a comprehensive performance evaluation of the designated strategies,
// applied to a specified collection of assets. In the absence of explicitly defined
// assets, encompasses all assets within the repository. Likewise, in the absence of
// explicitly defined strategies, encompasses all the registered strategies.
func (b *Backtest) Run() error {
	// When asset names are absent, considers all assets within the
	// provided repository for evaluation.
	if len(b.Names) == 0 {
		assets, err := b.repository.Assets()
		if err != nil {
			return err
		}

		b.Names = assets
	}

	// When strategies are absent, considers all strategies.
	if len(b.Strategies) == 0 {
		b.Strategies = b.allStrategies()
	}

	// Make sure that output directory exists.
	err := os.MkdirAll(b.outputDir, 0o700)
	if err != nil {
		return err
	}

	// Run the backtest workers and get the resutls.
	results := b.runWorkers()

	// Write the backtest report.
	return b.writeReport(results)
}

// allStrategies returns a slice containing references to all available strategies.
func (*Backtest) allStrategies() []Strategy {
	return []Strategy{
		NewBuyAndHoldStrategy(),
	}
}

// runWorkers initiates and manages the execution of multiple backtest workers.
func (b *Backtest) runWorkers() []*backtestResult {
	names := helper.SliceToChan(b.Names)
	results := make(chan *backtestResult)
	wg := &sync.WaitGroup{}

	for i := 0; i < b.Workers; i++ {
		wg.Add(1)
		go b.worker(names, results, wg)
	}

	// Wait for all workers to finish.
	go func() {
		wg.Wait()
		close(results)
	}()

	resultsSlice := helper.ChanToSlice(results)

	// Sort the backtest results by the outcomes.
	slices.SortFunc(resultsSlice, func(a, b *backtestResult) int {
		return int(b.Outcome - a.Outcome)
	})

	return resultsSlice
}

// strategyReportFileName
func (*Backtest) strategyReportFileName(assetName, strategyName string) string {
	return fmt.Sprintf("%s - %s.html", assetName, strategyName)
}

// worker is a backtesting worker that concurrently executes backtests for individual
// assets. It receives asset names from the provided channel, and performs backtests
// using the given strategies.
func (b *Backtest) worker(names <-chan string, bestResults chan<- *backtestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	since := time.Now().AddDate(0, 0, -b.LastDays)

	for name := range names {
		log.Printf("Backtesting %s...", name)
		snapshots, err := b.repository.GetSince(name, since)
		if err != nil {
			log.Printf("Unable to retrieve the snapshots for %s (%v)", name, err)
			continue
		}

		// We don't expect the snapshots to be a stream during backtesting.
		snapshotsSlice := helper.ChanToSlice(snapshots)

		results := make([]*backtestResult, 0, len(b.Strategies))

		for _, st := range b.Strategies {
			snapshotCopies := helper.Duplicate(helper.SliceToChan(snapshotsSlice), 2)

			actions, outcomes := ComputeWithOutcome(st, snapshotCopies[0])
			report := st.Report(snapshotCopies[1])

			actionsSplice := helper.Duplicate(actions, 3)

			actions = helper.Last(DenormalizeActions(actionsSplice[0]), 1)
			sinces := helper.Last(helper.Since[Action, int](actionsSplice[1]), 1)
			outcomes = helper.Last(outcomes, 1)
			transactions := helper.Last(CountTransactions(actionsSplice[2]), 1)

			err := report.WriteToFile(path.Join(b.outputDir, b.strategyReportFileName(name, st.Name())))
			if err != nil {
				log.Printf("Unable to write report for %s (%v)", name, err)
				continue
			}

			results = append(results, &backtestResult{
				AssetName:    name,
				StrategyName: st.Name(),
				Action:       <-actions,
				Since:        <-sinces,
				Outcome:      <-outcomes * 100,
				Transactions: <-transactions,
			})
		}

		// Sort the backtest results by the outcomes.
		slices.SortFunc(results, func(a, b *backtestResult) int {
			return int(b.Outcome - a.Outcome)
		})

		// Report the best result for the current asset.
		log.Printf("Best outcome for %s is %.2f%% with %s.", name, results[0].Outcome, results[0].StrategyName)
		bestResults <- results[0]

		// Write the asset report.
		err = b.writeAssetReport(name, results)
		if err != nil {
			log.Printf("Unable to write report with %v.", err)
		}
	}
}

// writeAssetReport generates a detailed report for the asset, summarizing the backtest results.
func (b *Backtest) writeAssetReport(name string, results []*backtestResult) error {
	type Model struct {
		AssetName string
		Results   []*backtestResult
	}

	model := Model{
		AssetName: name,
		Results:   results,
	}

	file, err := os.Create(filepath.Join(b.outputDir, fmt.Sprintf("%s.html", name)))
	if err != nil {
		return err
	}

	tmpl, err := template.New("report").Parse(backtestAssetReportTmpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, model)
	if err != nil {
		return err
	}

	return file.Close()
}

// writeReport generates a detailed report for the backtest results.
func (b *Backtest) writeReport(results []*backtestResult) error {
	type Model struct {
		Results []*backtestResult
	}

	model := Model{
		Results: results,
	}

	file, err := os.Create(filepath.Join(b.outputDir, "index.html"))
	if err != nil {
		return err
	}

	tmpl, err := template.New("report").Parse(backtestReportTmpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, model)
	if err != nil {
		return err
	}

	return file.Close()
}
