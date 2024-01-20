// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

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
	var tiingoKey string
	var targetBase string
	var minusDays int
	var workers int
	var delay int

	fmt.Fprintln(os.Stderr, "Indicator Sync")
	fmt.Fprintln(os.Stderr, "Copyright (c) 2021-2024 Onur Cinar.")
	fmt.Fprintln(os.Stderr, "The source code is provided under GNU AGPLv3 License.")
	fmt.Fprintln(os.Stderr, "https://github.com/cinar/indicator/v2")
	fmt.Fprintln(os.Stderr)

	flag.StringVar(&tiingoKey, "key", "", "tiingo service api key")
	flag.StringVar(&targetBase, "target", ".", "target repository base directory")
	flag.IntVar(&minusDays, "days", 0, "lookback period in days for the new assets")
	flag.IntVar(&workers, "workers", asset.DefaultSyncWorkers, "number of concurrent workers")
	flag.IntVar(&delay, "delay", asset.DefaultSyncDelay, "delay between each get")
	flag.Parse()

	if tiingoKey == "" {
		log.Fatal("Tiingo API key required")
	}

	defaultStartDate := time.Now().AddDate(0, 0, -minusDays)

	source := asset.NewTiingoRepository(tiingoKey)
	target := asset.NewFileSystemRepository(targetBase)

	sync := asset.NewSync()
	sync.Workers = workers
	sync.Delay = delay
	sync.Assets = flag.Args()

	err := sync.Run(source, target, defaultStartDate)
	if err != nil {
		log.Fatal(err)
	}
}
