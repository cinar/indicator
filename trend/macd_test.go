// Copyright (c) 2021-2024 Onur Cinar.
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
