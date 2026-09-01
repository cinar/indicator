// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator backtest command line program.
package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/strategy"
)

func main() {
	var repositoryName string
	var repositoryConfig string
	var reportName string
	var reportConfig string
	var strategyNames string
	var listStrategies bool
	var workers int
	var lastDays int
	var addSplits bool
	var addAnds bool

	stdErr := log.New(os.Stderr, "", 0)
	stdErr.Println("Indicator Backtest")
	stdErr.Println("Copyright (c) 2021-2026 The Indicator Authors.")
	stdErr.Println("The source code is provided under GNU AGPLv3 License.")
	stdErr.Println("https://github.com/cinar/indicator")
	stdErr.Println()
	stdErr.Println("DISCLAIMER: For educational and research purposes only. Not investment or financial advice.")
	stdErr.Println("Backtested performance is hypothetical and does not guarantee future results.")
	stdErr.Println()

	flag.StringVar(&repositoryName, "repository-name", "filesystem", "repository name")
	flag.StringVar(&repositoryConfig, "repository-config", "", "repository config")
	flag.StringVar(&reportName, "report-name", "html", "report name")
	flag.StringVar(&reportConfig, "report-config", ".", "report type")
	flag.StringVar(&strategyNames, "strategies", "", "comma-separated list of strategy names to backtest (see -list-strategies)")
	flag.BoolVar(&listStrategies, "list-strategies", false, "list the available strategy names and exit")
	flag.IntVar(&workers, "workers", backtest.DefaultBacktestWorkers, "number of concurrent workers")
	flag.IntVar(&lastDays, "last", backtest.DefaultLastDays, "number of days to do backtest")
	flag.BoolVar(&addSplits, "splits", false, "add the split strategies")
	flag.BoolVar(&addAnds, "ands", false, "add the and strategies")
	flag.Parse()

	logger := slog.Default()

	if listStrategies {
		for _, name := range StrategyNames() {
			stdErr.Println(name)
		}

		return
	}

	if strategyNames == "" {
		logger.Error("No strategies specified. Provide one or more with -strategies (comma-separated), or see -list-strategies for the available names.")
		os.Exit(1)
	}

	source, err := asset.NewRepository(repositoryName, repositoryConfig)
	if err != nil {
		logger.Error("Unable to initialize source.", "error", err)
		os.Exit(1)
	}

	report, err := backtest.NewReport(reportName, reportConfig)
	if err != nil {
		logger.Error("Unable to initialize report.", "error", err)
		os.Exit(1)
	}

	backtester := backtest.NewBacktest(source, report)
	backtester.Workers = workers
	backtester.LastDays = lastDays
	backtester.Logger = logger
	backtester.Names = append(backtester.Names, flag.Args()...)

	for _, name := range strings.Split(strategyNames, ",") {
		name = strings.TrimSpace(name)

		s, err := NewStrategy(name)
		if err != nil {
			logger.Error("Unable to initialize strategy.", "name", name, "error", err)
			os.Exit(1)
		}

		backtester.Strategies = append(backtester.Strategies, s)
	}

	if addSplits {
		backtester.Strategies = append(backtester.Strategies, strategy.AllSplitStrategies(backtester.Strategies)...)
	}

	if addAnds {
		backtester.Strategies = append(backtester.Strategies, strategy.AllAndStrategies(backtester.Strategies)...)
	}

	err = backtester.Run()
	if err != nil {
		logger.Error("Unable to run backtest.", "error", err)
		os.Exit(1)
	}
}
