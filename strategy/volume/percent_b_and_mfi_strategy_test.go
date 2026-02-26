// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy/volume"
)

func TestPercentBandMFIStrategy(t *testing.T) {
	pbms := volume.NewPercentBandMFIStrategy()

	snapshots := make(chan *asset.Snapshot, 100)
	for i := 0; i < 100; i++ {
		snapshots <- &asset.Snapshot{
			High:  float64(i + 10),
			Low:   float64(i),
			Close: float64(i + 5),
			Volume: 1000,
		}
	}
	close(snapshots)

	actions := pbms.Compute(snapshots)
	helper.Drain(actions)
}

func TestPercentBandMFIStrategyReport(t *testing.T) {
	pbms := volume.NewPercentBandMFIStrategy()

	snapshots := make(chan *asset.Snapshot, 100)
	for i := 0; i < 100; i++ {
		snapshots <- &asset.Snapshot{
			High:  float64(i + 10),
			Low:   float64(i),
			Close: float64(i + 5),
			Volume: 1000,
		}
	}
	close(snapshots)

	report := pbms.Report(snapshots)
	if report == nil {
		t.Fatal("expected report")
	}
}
