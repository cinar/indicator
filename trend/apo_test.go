// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestApo(t *testing.T) {
	type ApoData struct {
		Close float64
		Apo   float64
	}

	input, err := helper.ReadFromCsvFile[ApoData]("testdata/apo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(a *ApoData) float64 { return a.Close })
	expected := helper.Map(inputs[1], func(a *ApoData) float64 { return a.Apo })

	apo := trend.NewApo[float64]()
	actual := helper.RoundDigits(apo.Compute(closing), 2)
	expected = helper.Skip(expected, apo.SlowPeriod-1)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApoWithCustomSmoothing(t *testing.T) {
	type ApoData struct {
		Close float64
		Apo   float64
	}

	input, err := helper.ReadFromCsvFile[ApoData]("testdata/apo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 3)
	actualClosing := helper.Map(inputs[0], func(a *ApoData) float64 { return a.Close })
	fastClosing := helper.Map(inputs[1], func(a *ApoData) float64 { return a.Close })
	slowClosing := helper.Map(inputs[2], func(a *ApoData) float64 { return a.Close })

	apo := trend.NewApo[float64]()
	apo.FastSmoothing = 1.5
	apo.SlowSmoothing = 2.5
	actual := helper.RoundDigits(apo.Compute(actualClosing), 2)

	// Compute the expected APO independently, mirroring Apo's Fast - Slow
	// composition, to confirm the custom smoothing values are honored.
	fastEma := trend.NewEma[float64]()
	fastEma.Period = apo.FastPeriod
	fastEma.Smoothing = apo.FastSmoothing
	fast := helper.Skip(fastEma.Compute(fastClosing), apo.SlowPeriod-apo.FastPeriod)

	slowEma := trend.NewEma[float64]()
	slowEma.Period = apo.SlowPeriod
	slowEma.Smoothing = apo.SlowSmoothing
	slow := slowEma.Compute(slowClosing)

	expected := helper.RoundDigits(helper.Subtract(fast, slow), 2)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApoString(t *testing.T) {
	expected := "APO(14,30)"
	actual := trend.NewApo[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
