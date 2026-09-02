// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestTdSequential(t *testing.T) {
	type Data struct {
		Close         float64 `header:"Close"`
		BuySetup      float64 `header:"BuySetup"`
		SellSetup     float64 `header:"SellSetup"`
		BuyCountdown  float64 `header:"BuyCountdown"`
		SellCountdown float64 `header:"SellCountdown"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/td_sequential.csv")
	if err != nil {
		t.Fatal(err)
	}

	// Same pattern as RSI test
	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 {
		return d.BuySetup*1000 + d.SellSetup*100 + d.BuyCountdown*10 + d.SellCountdown
	})

	td := momentum.NewTdSequential[float64]()
	buySetup, sellSetup, buyCountdown, sellCountdown := td.Compute(closings)

	// Use Operate4 to combine all 4 outputs into one channel
	// Then use CheckEquals to compare with expected
	combined := helper.Operate4(buySetup, sellSetup, buyCountdown, sellCountdown,
		func(bs, ss, bc, sc float64) float64 {
			return bs*1000 + ss*100 + bc*10 + sc
		})

	// Skip to account for idle period
	combined = helper.Skip(combined, td.IdlePeriod())
	expected = helper.Skip(expected, td.IdlePeriod())

	err = helper.CheckEquals(combined, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// tdSequentialRow bundles one bar's worth of TD Sequential outputs so that
// all four output channels can be drained together without deadlocking.
type tdSequentialRow struct {
	buySetup      float64
	sellSetup     float64
	buyCountdown  float64
	sellCountdown float64
}

func tdSequentialRows(closings []float64) []tdSequentialRow {
	td := momentum.NewTdSequential[float64]()
	buySetup, sellSetup, buyCountdown, sellCountdown := td.Compute(helper.SliceToChan(closings))

	combined := helper.Operate4(buySetup, sellSetup, buyCountdown, sellCountdown,
		func(bs, ss, bc, sc float64) tdSequentialRow {
			return tdSequentialRow{
				buySetup:      bs,
				sellSetup:     ss,
				buyCountdown:  bc,
				sellCountdown: sc,
			}
		})

	return helper.ChanToSlice(combined)
}

// TestTdSequentialSetupCapsAtSetupPeriod checks that a long streak of
// consecutive lower closes does not push the buy setup count past the
// documented SetupPeriod (9). Without the fix, currentBuySetup keeps
// climbing (10, 11, 12, ...) for as long as the streak lasts.
func TestTdSequentialSetupCapsAtSetupPeriod(t *testing.T) {
	// 20 strictly decreasing closes. With the default Lookback of 4, every
	// close from index 4 onward is lower than the close 4 bars back, so the
	// buy setup count would climb every bar: 1, 2, 3, ... 16 by the last bar
	// if left uncapped. It must instead reach 9 (at index 12) and hold.
	closings := make([]float64, 0, 20)
	price := 100.0

	for i := 0; i < 20; i++ {
		closings = append(closings, price)
		price--
	}

	rows := tdSequentialRows(closings)

	for i, row := range rows {
		if row.buySetup > momentum.DefaultTdSequentialSetupPeriod {
			t.Fatalf("row %d: buySetup = %v, want <= %v", i, row.buySetup, momentum.DefaultTdSequentialSetupPeriod)
		}
	}

	// Sanity check that the streak actually reaches the cap and holds there,
	// rather than the assertion above trivially passing on a series that
	// never gets close to 9.
	if rows[12].buySetup != 9 {
		t.Fatalf("row 12: buySetup = %v, want 9", rows[12].buySetup)
	}

	if rows[19].buySetup != 9 {
		t.Fatalf("row 19: buySetup = %v, want 9 (clamped, not 16)", rows[19].buySetup)
	}
}

// TestTdSequentialCountdownCrossCancel checks that an in-progress buy
// countdown is cancelled the moment an opposite (sell) setup completes, as
// required by the TD Sequential methodology. Without the fix, buyCountdown
// keeps accumulating independently of the newly completed sell setup.
func TestTdSequentialCountdownCrossCancel(t *testing.T) {
	closings := []float64{
		100, 99, 98, 97, // fill the Lookback(4) buffer
		96, 95, 94, 93, 92, 91, 90, 89, // buy setup counts 1..8
		88,         // buy setup completes at 9, buy countdown starts
		87, 86, 85, // buy countdown advances: 1, 2, 3, 4
		90, 95, 100, 105, 110, 115, 120, 125, // reverse hard: sell setup counts -1..-8
		130, // sell setup completes at -9: must cancel the buy countdown
		120, // pullback that WOULD advance the buy countdown to 5 if it had
		// not been cancelled (close <= close 2 bars ago, i.e. 120 <= 125)
	}

	rows := tdSequentialRows(closings)

	// Confirm the series actually builds up a buy countdown before the
	// reversal, so the later reset is meaningful and not a no-op.
	if rows[15].buyCountdown != 4 {
		t.Fatalf("row 15: buyCountdown = %v, want 4 (in progress before reversal)", rows[15].buyCountdown)
	}

	// Confirm the reversal actually completes a sell setup.
	if rows[24].sellSetup != -9 {
		t.Fatalf("row 24: sellSetup = %v, want -9 (setup completed)", rows[24].sellSetup)
	}

	// The completed sell setup must cancel the in-progress buy countdown.
	if rows[24].buyCountdown != 0 {
		t.Fatalf("row 24: buyCountdown = %v, want 0 (cancelled by opposite setup)", rows[24].buyCountdown)
	}

	// The cancellation must stick: the following pullback bar satisfies the
	// buy countdown's own comparison (120 <= close 2 bars ago = 125), so
	// without cross-cancellation the count would climb back to 5.
	if rows[25].buyCountdown != 0 {
		t.Fatalf("row 25: buyCountdown = %v, want 0 (stays cancelled, not 5)", rows[25].buyCountdown)
	}
}

func TestTdSequentialString(t *testing.T) {
	expected := "TDSEQUENTIAL(4,2,9,13)"
	actual := momentum.NewTdSequential[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
