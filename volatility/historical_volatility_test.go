// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func TestHistoricalVolatilityDefaultAndString(t *testing.T) {
	hv := volatility.NewHistoricalVolatility[float64]()

	if hv.Period != volatility.DefaultHistoricalVolatilityPeriod {
		t.Fatalf("expected period %d, got %d", volatility.DefaultHistoricalVolatilityPeriod, hv.Period)
	}

	if hv.IdlePeriod() != volatility.DefaultHistoricalVolatilityPeriod {
		t.Fatalf("expected idle period %d, got %d", volatility.DefaultHistoricalVolatilityPeriod, hv.IdlePeriod())
	}

	if hv.String() != "HV(21)" {
		t.Fatalf("expected HV(21), got %s", hv.String())
	}
}

func TestHistoricalVolatilityWithInvalidPeriod(t *testing.T) {
	hv := volatility.NewHistoricalVolatilityWithPeriod[float64](0)

	if hv.Period != volatility.DefaultHistoricalVolatilityPeriod {
		t.Fatalf("expected default period %d, got %d", volatility.DefaultHistoricalVolatilityPeriod, hv.Period)
	}
}

func TestHistoricalVolatility(t *testing.T) {
	type Data struct {
		Close float64 `header:"Close"`
		Hv    float64 `header:"Hv"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/historical_volatility.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Hv })

	hv := volatility.NewHistoricalVolatility[float64]()
	actual := hv.Compute(closings)
	actual = helper.RoundDigits(actual, 8)

	expected = helper.Skip(expected, hv.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

