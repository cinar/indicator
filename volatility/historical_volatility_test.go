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

func TestHistoricalVolatilityCompute(t *testing.T) {
	prices := helper.SliceToChan([]float64{100, 110, 121, 133.1, 146.41, 161.051})

	hv := volatility.NewHistoricalVolatilityWithPeriod[float64](2)
	actuals := hv.Compute(prices)
	actuals = helper.RoundDigits(actuals, 8)

	expected := helper.SliceToChan([]float64{0, 0, 0})
	err := helper.CheckEquals(actuals, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHistoricalVolatilityPreviousPriceZero(t *testing.T) {
	prices := helper.SliceToChan([]float64{0, 5, 10, 20, 40})

	hv := volatility.NewHistoricalVolatilityWithPeriod[float64](2)
	actuals := hv.Compute(prices)
	actuals = helper.RoundDigits(actuals, 8)

	expected := helper.SliceToChan([]float64{0.5, 0})
	err := helper.CheckEquals(actuals, expected)
	if err != nil {
		t.Fatal(err)
	}
}
