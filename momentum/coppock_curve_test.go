// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestCoppockCurveToStringAndIdlePeriod(t *testing.T) {
	cc := momentum.NewCoppockCurve[float64]()
	if cc.String() != "CoppockCurve(14,11,10)" {
		t.Fatalf("unexpected string: %s", cc.String())
	}
	if cc.IdlePeriod() != 23 {
		t.Fatalf("unexpected IdlePeriod: %d", cc.IdlePeriod())
	}
}

func TestCoppockCurveWithDefaultsAndZero(t *testing.T) {
	cc := momentum.NewCoppockCurveWithPeriods[float64](0, 0, 0)
	if cc.String() != "CoppockCurve(14,11,10)" {
		t.Fatalf("unexpected string: %s", cc.String())
	}
}

func TestCoppockCurveRoc2Longer(t *testing.T) {
	// Roc1: 2, Roc2: 3, Wma: 2
	// IdlePeriod: max(2, 3) + 2 - 1 = 3 + 1 = 4
	cc := momentum.NewCoppockCurveWithPeriods[float64](2, 3, 2)
	if cc.IdlePeriod() != 4 {
		t.Fatalf("unexpected IdlePeriod: %d", cc.IdlePeriod())
	}

	closings := helper.SliceToChan([]float64{10, 12, 11, 13, 15, 14, 16})
	actuals := cc.Compute(closings)
	
	count := 0
	for range actuals {
		count++
	}
	
	// Total 7, skip 4 = 3
	if count != 3 {
		t.Fatalf("expected 3 values, got %d", count)
	}
}

func TestCoppockCurveTestdata(t *testing.T) {
	type Data struct {
		Close        float64 `header:"Close"`
		CoppockCurve float64 `header:"CoppockCurve"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/coppock_curve.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.CoppockCurve })

	cc := momentum.NewCoppockCurve[float64]()
	actuals := cc.Compute(closings)
	actuals = helper.RoundDigits(actuals, 2)

	expected = helper.Skip(expected, cc.IdlePeriod())

	err = helper.CheckEquals(actuals, expected)
	if err != nil {
		t.Fatal(err)
	}
}
