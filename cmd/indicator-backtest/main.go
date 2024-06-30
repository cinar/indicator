// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator backtest command line program.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/compound"
	"github.com/cinar/indicator/v2/strategy/momentum"
	"github.com/cinar/indicator/v2/strategy/trend"
	"github.com/cinar/indicator/v2/strategy/volatility"
)

func main() {
	var repositoryDir string
	var outputDir string
	var workers int
	var lastDays int
	var writeStrategyRerpots bool
	var addSplits bool
	var addAnds bool
	var dateFormat string

	fmt.Fprintln(os.Stderr, "Indicator Backtest")
	fmt.Fprintln(os.Stderr, "Copyright (c) 2021-2024 Onur Cinar.")
	fmt.Fprintln(os.Stderr, "The source code is provided under GNU AGPLv3 License.")
	fmt.Fprintln(os.Stderr, "https://github.com/cinar/indicator")
	fmt.Fprintln(os.Stderr)

	flag.StringVar(&repositoryDir, "repository", ".", "file system repository directory")
	flag.StringVar(&outputDir, "output", ".", "output directory")
	flag.IntVar(&workers, "workers", strategy.DefaultBacktestWorkers, "number of concurrent workers")
	flag.IntVar(&lastDays, "last", strategy.DefaultLastDays, "number of days to do backtest")
	flag.BoolVar(&writeStrategyRerpots, "write-strategy-reports", strategy.DefaultWriteStrategyReports, "write individual strategy reports")
	flag.BoolVar(&addSplits, "splits", false, "add the split strategies")
	flag.BoolVar(&addAnds, "ands", false, "add the and strategies")
	flag.StringVar(&dateFormat, "date-format", helper.DefaultReportDateFormat, "date format to use")
	flag.Parse()

	flag.Parse()

	repository := asset.NewFileSystemRepository(repositoryDir)

	backtest := strategy.NewBacktest(repository, outputDir)
	backtest.Workers = workers
	backtest.LastDays = lastDays
	backtest.WriteStrategyReports = writeStrategyRerpots
	backtest.DateFormat = dateFormat
	backtest.Names = append(backtest.Names, flag.Args()...)
	backtest.Strategies = append(backtest.Strategies, compound.AllStrategies()...)
	backtest.Strategies = append(backtest.Strategies, momentum.AllStrategies()...)
	backtest.Strategies = append(backtest.Strategies, strategy.AllStrategies()...)
	backtest.Strategies = append(backtest.Strategies, trend.AllStrategies()...)
	backtest.Strategies = append(backtest.Strategies, volatility.AllStrategies()...)

	if addSplits {
		backtest.Strategies = append(backtest.Strategies, strategy.AllSplitStrategies(backtest.Strategies)...)
	}

	if addAnds {
		backtest.Strategies = append(backtest.Strategies, strategy.AllAndStrategies(backtest.Strategies)...)
	}

	err := backtest.Run()
	if err != nil {
		log.Fatal(err)
	}
}
