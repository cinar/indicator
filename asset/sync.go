// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultSyncWorkers is the default number of workers to use to synchronize.
	DefaultSyncWorkers = 1

	// DefaultSyncDelay is the default delay in seconds between each get request.
	DefaultSyncDelay = 5
)

// Sync represents the configuration parameters for synchronizing assets between repositories.
type Sync struct {
	// Number of workers to use.
	Workers int

	// Delay between repository get requests to minimize the load to the remote server.
	Delay int

	// Assets is the name of the assets to be synced. If it is empty, all assets in the target repository
	// will be synced instead.
	Assets []string
}

// NewSync function initializes a new sync instance with the default parameters.
func NewSync() *Sync {
	return &Sync{
		Workers: DefaultSyncWorkers,
		Delay:   DefaultSyncDelay,
		Assets:  []string{},
	}
}

// Run synchronizes assets between the source and target repositories using multi-worker concurrency.
func (s *Sync) Run(source, target Repository, defaultStartDate time.Time) error {
	if len(s.Assets) == 0 {
		log.Print("No asset names provided. Syncing in all assets in the target repository.")

		assets, err := target.Assets()
		if err != nil {
			return err
		}

		s.Assets = assets
	}

	log.Printf("Will sync %d assets.", len(s.Assets))
	jobs := helper.SliceToChan(s.Assets)

	hasErrors := false
	wg := &sync.WaitGroup{}

	for i := 0; i < s.Workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for name := range jobs {
				lastDate, err := target.LastDate(name)
				if err == nil {
					lastDate = lastDate.AddDate(0, 0, 1)
				} else {
					lastDate = defaultStartDate
				}

				log.Printf("Syncing %s starting %s", name, lastDate.Format("2006-01-02"))

				snapshots, err := source.GetSince(name, lastDate)
				if err != nil {
					log.Printf("GetSince failed for %s with %v", name, err)
					hasErrors = true
					continue
				}

				err = target.Append(name, snapshots)
				if err != nil {
					log.Printf("Append failed for %s with %v", name, err)
					hasErrors = true
					continue
				}

				time.Sleep(time.Duration(s.Delay) * time.Second)
			}
		}()
	}

	wg.Wait()

	if hasErrors {
		return errors.New("has errors")
	}

	return nil
}
