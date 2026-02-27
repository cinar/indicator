// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestChop(t *testing.T) {
	type Data struct {
		High     float64 `header:"High"`
		Low      float64 `header:"Low"`
		Close    float64 `header:"Close"`
		Expected float64 `header:"Expected"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/chop.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)

	highs := helper.Map(inputs[0], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[1], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[2], func(d *Data) float64 { return d.Close })

	chop := volatility.NewChop[float64]()
	actuals := chop.Compute(highs, lows, closings)
	actuals = helper.RoundDigits(actuals, 6)

	inputs[3] = helper.Skip(inputs[3], chop.IdlePeriod())

	for data := range inputs[3] {
		actual, ok := <-actuals
		if !ok {
			t.Fatal("actuals channel closed early")
		}
		if actual != data.Expected {
			t.Fatalf("actual %v expected %v for High=%v Low=%v Close=%v", actual, data.Expected, data.High, data.Low, data.Close)
		}
	}

	if _, ok := <-actuals; ok {
		t.Fatal("actuals channel should be closed")
	}

	if chop.String() != "CHOP(14)" {
		t.Fatalf("expected CHOP(14) but got %s", chop.String())
	}

	if chop.IdlePeriod() != 14 {
		t.Fatalf("expected 14 but got %d", chop.IdlePeriod())
	}
}

func TestNewChop(t *testing.T) {
	chop := volatility.NewChop[float64]()
	if chop == nil {
		t.Fatal("NewChop should not be nil")
	}
}

func TestNewChopWithPeriod(t *testing.T) {
	chop := volatility.NewChopWithPeriod[float64](10)
	if chop.Period != 10 {
		t.Fatalf("expected 10 but got %d", chop.Period)
	}
}

func TestChopDiffZero(t *testing.T) {
	highs := helper.SliceToChan([]float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10})
	lows := helper.SliceToChan([]float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10})
	closings := helper.SliceToChan([]float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10})

	chop := volatility.NewChopWithPeriod[float64](2)
	actuals := chop.Compute(highs, lows, closings)

	for actual := range actuals {
		if actual != 0 {
			t.Fatalf("expected 0 but got %v", actual)
		}
	}
}
