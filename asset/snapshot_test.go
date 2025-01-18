// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"testing"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

func TestSnapshotsAs(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	snapshotsCopies := helper.Duplicate(snapshots, 7)

	dates := asset.SnapshotsAsDates(snapshotsCopies[1])
	openings := asset.SnapshotsAsOpenings(snapshotsCopies[2])
	highs := asset.SnapshotsAsHighs(snapshotsCopies[3])
	lows := asset.SnapshotsAsLows(snapshotsCopies[4])
	closings := asset.SnapshotsAsClosings(snapshotsCopies[5])
	volumes := asset.SnapshotsAsVolumes(snapshotsCopies[6])

	for snapshot := range snapshotsCopies[0] {
		date := <-dates
		opening := <-openings
		high := <-highs
		low := <-lows
		closing := <-closings
		volume := <-volumes

		if !date.Equal(snapshot.Date) {
			t.Fatalf("actual %v expected %v", date, snapshot.Date)
		}

		if opening != snapshot.Open {
			t.Fatalf("actual %v expected %v", opening, snapshot.Open)
		}

		if high != snapshot.High {
			t.Fatalf("actual %v expected %v", high, snapshot.High)
		}

		if low != snapshot.Low {
			t.Fatalf("actual %v expected %v", low, snapshot.Low)
		}

		if closing != snapshot.Close {
			t.Fatalf("actual %v expected %v", closing, snapshot.Close)
		}

		if volume != snapshot.Volume {
			t.Fatalf("actual %v expected %v", volume, snapshot.Volume)
		}
	}
}
