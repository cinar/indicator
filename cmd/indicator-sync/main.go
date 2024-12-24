// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator sync command line program.
package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/cinar/indicator/v2/asset"
)

func main() {
	var sourceName string
	var sourceConfig string
	var targetName string
	var targetConfig string
	var minusDays int
	var workers int
	var delay int

	stdErr := log.New(os.Stderr, "", 0)
	stdErr.Println("Indicator Sync")
	stdErr.Println("Copyright (c) 2021-2024 Onur Cinar.")
	stdErr.Println("The source code is provided under GNU AGPLv3 License.")
	stdErr.Println("https://github.com/cinar/indicator")
	stdErr.Println()

	flag.StringVar(&sourceName, "source-name", "tiingo", "source repository type")
	flag.StringVar(&sourceConfig, "source-config", "", "source repository config")
	flag.StringVar(&targetName, "target-name", "filesystem", "target repository type")
	flag.StringVar(&targetConfig, "target-config", "", "target repository config")
	flag.IntVar(&minusDays, "days", 0, "lookback period in days for the new assets")
	flag.IntVar(&workers, "workers", asset.DefaultSyncWorkers, "number of concurrent workers")
	flag.IntVar(&delay, "delay", asset.DefaultSyncDelay, "delay between each get")
	flag.Parse()

	logger := slog.Default()

	source, err := asset.NewRepository(sourceName, sourceConfig)
	if err != nil {
		logger.Error("Unable to initialize source.", "error", err)
		os.Exit(1)
	}

	target, err := asset.NewRepository(targetName, targetConfig)
	if err != nil {
		logger.Error("Unable to initialize target.", "error", err)
		os.Exit(1)
	}

	defaultStartDate := time.Now().AddDate(0, 0, -minusDays)

	assets := flag.Args()
	if len(assets) == 0 {
		assets, err = source.Assets()
		if err != nil {
			logger.Error("Unable to get assets.", "error", err)
			os.Exit(1)
		}
	}

	sync := asset.NewSync()
	sync.Workers = workers
	sync.Delay = delay
	sync.Assets = assets
	sync.Logger = logger

	err = sync.Run(source, target, defaultStartDate)
	if err != nil {
		logger.Error("Unable to sync repositories.", "error", err)
		os.Exit(1)
	}
}
