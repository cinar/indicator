// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestAnnualizedHistoricalVolatilityDefaultAndString(t *testing.T) {
	ahv := volatility.NewAnnualizedHistoricalVolatility[float64]()

	if ahv.Hv.Period != volatility.DefaultAnnualizedHistoricalVolatilityPeriod {
		t.Fatalf("expected period %d, got %d", volatility.DefaultAnnualizedHistoricalVolatilityPeriod, ahv.Hv.Period)
	}

	if ahv.TradingDaysPerYear != volatility.DefaultTradingDaysPerYear {
		t.Fatalf("expected trading days %d, got %d", volatility.DefaultTradingDaysPerYear, ahv.TradingDaysPerYear)
	}

	if ahv.IdlePeriod() != volatility.DefaultAnnualizedHistoricalVolatilityPeriod {
		t.Fatalf("expected idle period %d, got %d", volatility.DefaultAnnualizedHistoricalVolatilityPeriod, ahv.IdlePeriod())
	}

	if ahv.String() != "AHV(21)" {
		t.Fatalf("expected AHV(21), got %s", ahv.String())
	}
}

func TestAnnualizedHistoricalVolatilityWithInvalidPeriod(t *testing.T) {
	ahv := volatility.NewAnnualizedHistoricalVolatilityWithPeriod[float64](0)

	if ahv.Hv.Period != volatility.DefaultAnnualizedHistoricalVolatilityPeriod {
		t.Fatalf("expected default period %d, got %d", volatility.DefaultAnnualizedHistoricalVolatilityPeriod, ahv.Hv.Period)
	}
}

func TestAnnualizedHistoricalVolatility(t *testing.T) {
	type Data struct {
		Close float64 `header:"Close"`
		Ahv   float64 `header:"Ahv"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/annualized_historical_volatility.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Ahv })

	ahv := volatility.NewAnnualizedHistoricalVolatility[float64]()
	actual := ahv.Compute(closings)
	actual = helper.RoundDigits(actual, 8)

	expected = helper.Skip(expected, ahv.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
