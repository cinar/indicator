// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator sync command line program.
package main

import (
	"flag"
	"fmt"
	"log"
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

	fmt.Fprintln(os.Stderr, "Indicator Sync")
	fmt.Fprintln(os.Stderr, "Copyright (c) 2021-2024 Onur Cinar.")
	fmt.Fprintln(os.Stderr, "The source code is provided under GNU AGPLv3 License.")
	fmt.Fprintln(os.Stderr, "https://github.com/cinar/indicator")
	fmt.Fprintln(os.Stderr)

	flag.StringVar(&sourceName, "source-name", "tiingo", "source repository type")
	flag.StringVar(&sourceConfig, "source-config", "", "source repository config")
	flag.StringVar(&targetName, "target-name", "filesystem", "target repository type")
	flag.StringVar(&targetConfig, "target-config", "", "target repository config")
	flag.IntVar(&minusDays, "days", 0, "lookback period in days for the new assets")
	flag.IntVar(&workers, "workers", asset.DefaultSyncWorkers, "number of concurrent workers")
	flag.IntVar(&delay, "delay", asset.DefaultSyncDelay, "delay between each get")
	flag.Parse()

	source, err := asset.NewRepository(sourceName, sourceConfig)
	if err != nil {
		log.Fatalf("unable to initialize source: %v", err)
	}

	target, err := asset.NewRepository(targetName, targetConfig)
	if err != nil {
		log.Fatalf("unable to initialize target: %v", err)
	}

	defaultStartDate := time.Now().AddDate(0, 0, -minusDays)

	assets := flag.Args()
	if len(assets) == 0 {
		assets, err = source.Assets()
		if err != nil {
			log.Fatalf("unable to get assets: %v", err)
		}
	}

	sync := asset.NewSync()
	sync.Workers = workers
	sync.Delay = delay
	sync.Assets = assets

	err = sync.Run(source, target, defaultStartDate)
	if err != nil {
		log.Fatalf("unable to sync repositories: %v", err)
	}
}
