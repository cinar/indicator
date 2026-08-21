// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/registry"
	"github.com/cinar/indicator/v2/trend"
)

// testSnapshots builds a deterministic, non-flat series of snapshots so
// that indicators with a denominator (e.g. CCI) don't see an all-zero
// window.
func testSnapshots(count int) []*asset.Snapshot {
	snapshots := make([]*asset.Snapshot, count)

	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	for i := range snapshots {
		closing := 100 + 10*math.Sin(float64(i)/3)

		snapshots[i] = &asset.Snapshot{
			Date:   date.AddDate(0, 0, i),
			Open:   closing - 1,
			High:   closing + 2,
			Low:    closing - 2,
			Close:  closing,
			Volume: 1000 + float64(i),
		}
	}

	return snapshots
}

func TestNames(t *testing.T) {
	names := registry.Names()

	for _, want := range []string{"sma", "ema", "macd", "cci", "rsi", "bollinger_bands", "obv", "vwap"} {
		found := false

		for _, name := range names {
			if name == want {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected %q to be registered, got %v", want, names)
		}
	}
}

func TestGetUnknown(t *testing.T) {
	if _, ok := registry.Get("does-not-exist"); ok {
		t.Fatal("expected unknown indicator to not be found")
	}
}

func TestRunUnknownIndicator(t *testing.T) {
	snapshots := helper.SliceToChan(testSnapshots(10))

	_, err := registry.Run(context.Background(), "does-not-exist", nil, snapshots)
	if err == nil {
		t.Fatal("expected an error for an unknown indicator")
	}
}

// drain reads result.Dates and every result.Series channel in lockstep
// (one date, then one value per series, repeated) since they come from a
// shared, unbuffered fan-out and must be consumed together to avoid
// deadlocking the pipeline. It returns the dates and, for each series in
// order, its collected values.
func drain(t *testing.T, result *registry.Result) ([]time.Time, [][]float64) {
	t.Helper()

	var dates []time.Time
	values := make([][]float64, len(result.Series))

	for date := range result.Dates {
		dates = append(dates, date)

		for i, series := range result.Series {
			value, ok := <-series.Values
			if !ok {
				t.Fatalf("series %q ended before dates", series.Name)
			}

			values[i] = append(values[i], value)
		}
	}

	return dates, values
}

func TestRunSma(t *testing.T) {
	const period = 5

	data := testSnapshots(30)
	snapshots := helper.SliceToChan(data)

	result, err := registry.Run(context.Background(), "sma", map[string]string{"Period": "5"}, snapshots)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Series) != 1 || result.Series[0].Name != "sma" {
		t.Fatalf("unexpected series: %+v", result.Series)
	}

	dates, values := drain(t, result)
	actual := values[0]

	expectedSma := trend.NewSmaWithPeriod[float64](period)
	closings := asset.SnapshotsAsClosings(helper.SliceToChan(data))
	expected := helper.ChanToSlice(expectedSma.ComputeWithContext(context.Background(), closings))

	if len(dates) != len(data)-expectedSma.IdlePeriod() {
		t.Fatalf("expected %d dates after skipping the idle period, got %d", len(data)-expectedSma.IdlePeriod(), len(dates))
	}

	if len(actual) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("value %d: expected %v, got %v", i, expected[i], actual[i])
		}
	}

	if !dates[0].Equal(data[expectedSma.IdlePeriod()].Date) {
		t.Fatalf("expected first date to be %v, got %v", data[expectedSma.IdlePeriod()].Date, dates[0])
	}
}

func TestRunMacdWithNestedParam(t *testing.T) {
	snapshots := helper.SliceToChan(testSnapshots(60))

	result, err := registry.Run(context.Background(), "macd", map[string]string{
		"Ema1.Period": "6",
		"Ema2.Period": "13",
		"Ema3.Period": "5",
	}, snapshots)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Series) != 2 || result.Series[0].Name != "macd" || result.Series[1].Name != "signal" {
		t.Fatalf("unexpected series: %+v", result.Series)
	}

	dates, values := drain(t, result)

	if len(dates) == 0 {
		t.Fatal("expected at least one date")
	}

	if len(values[0]) != len(dates) || len(values[1]) != len(dates) {
		t.Fatalf("dates (%d), macd (%d), and signal (%d) are not aligned", len(dates), len(values[0]), len(values[1]))
	}
}

func TestRunCciMultiInput(t *testing.T) {
	data := testSnapshots(40)
	snapshots := helper.SliceToChan(data)

	result, err := registry.Run(context.Background(), "cci", map[string]string{"Period": "10"}, snapshots)
	if err != nil {
		t.Fatal(err)
	}

	dates, values := drain(t, result)

	if len(dates) != len(values[0]) {
		t.Fatalf("dates (%d) and values (%d) are not aligned", len(dates), len(values[0]))
	}

	const period = 10
	idlePeriod := period*2 - 2

	if len(dates) != len(data)-idlePeriod {
		t.Fatalf("expected %d dates, got %d", len(data)-idlePeriod, len(dates))
	}
}

func TestRunEveryRegisteredIndicator(t *testing.T) {
	// Some indicators (e.g. ConnorsRsi's 100-period PercentRank) need a
	// long lookback before they yield a single value.
	data := testSnapshots(400)

	for _, name := range registry.Names() {
		t.Run(name, func(t *testing.T) {
			snapshots := helper.SliceToChan(data)

			result, err := registry.Run(context.Background(), name, nil, snapshots)
			if err != nil {
				t.Fatal(err)
			}

			if len(result.Series) == 0 {
				t.Fatal("expected at least one output series")
			}

			dates, values := drain(t, result)

			if len(dates) == 0 {
				t.Fatal("expected at least one date after the idle period")
			}

			for i, series := range values {
				if len(series) != len(dates) {
					t.Fatalf("series %q has %d values, expected %d", result.Series[i].Name, len(series), len(dates))
				}
			}
		})
	}
}

func TestSetParamErrors(t *testing.T) {
	sma := trend.NewSma[float64]()

	if err := registry.SetParam(sma, "Period", "not-a-number"); err == nil {
		t.Fatal("expected an error for a non-numeric value")
	}

	if err := registry.SetParam(sma, "DoesNotExist", "1"); err == nil {
		t.Fatal("expected an error for an unknown param")
	}

	if err := registry.SetParam(sma, "Period", "10"); err != nil {
		t.Fatal(err)
	}

	if sma.Period != 10 {
		t.Fatalf("expected Period to be 10, got %d", sma.Period)
	}
}
