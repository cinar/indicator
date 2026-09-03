// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMacd(t *testing.T) {
	type Data struct {
		Close  float64
		Macd   float64
		Signal float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/macd.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	macd := trend.NewMacd[float64]()
	actualMacds, actualSignals := macd.Compute(closing)

	actualMacds = helper.RoundDigits(actualMacds, 2)
	actualSignals = helper.RoundDigits(actualSignals, 2)

	inputs[1] = helper.Skip(inputs[1], macd.IdlePeriod())

	for data := range inputs[1] {
		actualMacd := <-actualMacds
		actualSignal := <-actualSignals

		if actualMacd != data.Macd {
			t.Fatalf("actual %v expected %v", actualMacd, data.Macd)
		}

		if actualSignal != data.Signal {
			t.Fatalf("actual %v expected %v", actualSignal, data.Signal)
		}
	}
}

func TestMacdString(t *testing.T) {
	expected := "MACD(12,26,9)"
	actual := trend.NewMacd[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

// TestMacdWithReversedPeriods asserts that constructing a Macd with the
// fast and slow periods passed in reversed order (period1 > period2) is
// automatically normalized, and behaves identically to constructing it
// with the periods in the correct order.
func TestMacdWithReversedPeriods(t *testing.T) {
	type Data struct {
		Close  float64
		Macd   float64
		Signal float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/macd.csv")
	if err != nil {
		t.Fatal(err)
	}

	closing := helper.Map(input, func(d *Data) float64 { return d.Close })

	// Passed in reversed order: slow period first, fast period second.
	reversed := trend.NewMacdWithPeriod[float64](
		trend.DefaultMacdPeriod2,
		trend.DefaultMacdPeriod1,
		trend.DefaultMacdPeriod3,
	)

	expected := "MACD(12,26,9)"
	actual := reversed.String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}

	actualMacds, actualSignals := reversed.Compute(closing)

	actualMacds = helper.RoundDigits(actualMacds, 2)
	actualSignals = helper.RoundDigits(actualSignals, 2)

	inputs, err := helper.ReadFromCsvFile[Data]("testdata/macd.csv")
	if err != nil {
		t.Fatal(err)
	}

	expectedData := helper.Skip(inputs, reversed.IdlePeriod())

	for data := range expectedData {
		actualMacd := <-actualMacds
		actualSignal := <-actualSignals

		if actualMacd != data.Macd {
			t.Fatalf("actual %v expected %v", actualMacd, data.Macd)
		}

		if actualSignal != data.Signal {
			t.Fatalf("actual %v expected %v", actualSignal, data.Signal)
		}
	}
}
