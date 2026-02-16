// Copyright (c) 2021-2024 Onur Cinar.
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
