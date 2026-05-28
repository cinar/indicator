// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestPivotPoint(t *testing.T) {
	type Data struct {
		Open  float64 `header:"Open"`
		High  float64 `header:"High"`
		Low   float64 `header:"Low"`
		Close float64 `header:"Close"`
	}

	input, err := helper.ReadFromCsvFile[Data]("../helper/testdata/report.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	opens := helper.Map(inputs[0], func(d *Data) float64 { return d.Open })
	highs := helper.Map(inputs[1], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[2], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[3], func(d *Data) float64 { return d.Close })

	pp := trend.NewPivotPoint[float64]()
	results := pp.Compute(opens, highs, lows, closings)

	// First bar: 315.130005, 318.600006, 308.700012, 318.600006
	// Second bar Open: 319
	// Standard Pivot Point for second bar:
	// P = (318.600006 + 308.700012 + 318.600006) / 3 = 315.300008
	// R1 = 2 * 315.300008 - 308.700012 = 321.900004
	// S1 = 2 * 315.300008 - 318.600006 = 312.00001

	res := <-results

	if helper.RoundDigit(res.P, 6) != 315.300008 {
		t.Fatalf("expected P 315.300008, got %v", res.P)
	}

	if helper.RoundDigit(res.R1, 6) != 321.900004 {
		t.Fatalf("expected R1 321.900004, got %v", res.R1)
	}

	if helper.RoundDigit(res.S1, 6) != 312.00001 {
		t.Fatalf("expected S1 312.00001, got %v", res.S1)
	}
}

func TestPivotPointWoodie(t *testing.T) {
	type Data struct {
		Open  float64 `header:"Open"`
		High  float64 `header:"High"`
		Low   float64 `header:"Low"`
		Close float64 `header:"Close"`
	}

	input, err := helper.ReadFromCsvFile[Data]("../helper/testdata/report.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	opens := helper.Map(inputs[0], func(d *Data) float64 { return d.Open })
	highs := helper.Map(inputs[1], func(d *Data) float64 { return d.High })
	lows := helper.Map(inputs[2], func(d *Data) float64 { return d.Low })
	closings := helper.Map(inputs[3], func(d *Data) float64 { return d.Close })

	pp := trend.NewPivotPointWithMethod[float64](trend.PivotPointWoodie)
	results := pp.Compute(opens, highs, lows, closings)

	// First bar: H=318.600006, L=308.700012
	// Second bar: O=319
	// Woodie P = (318.600006 + 308.700012 + 2 * 319) / 4 = 316.3250045

	res := <-results

	if helper.RoundDigit(res.P, 7) != 316.3250045 {
		t.Fatalf("expected P 316.3250045, got %v", res.P)
	}
}

func TestPivotPointIdlePeriod(t *testing.T) {
	pp := trend.NewPivotPoint[float64]()
	if pp.IdlePeriod() != 1 {
		t.Fatalf("expected IdlePeriod 1, got %d", pp.IdlePeriod())
	}
}

func TestPivotPointString(t *testing.T) {
	pp := trend.NewPivotPoint[float64]()
	if pp.String() != "PivotPoint(Standard)" {
		t.Fatalf("expected PivotPoint(Standard), got %s", pp.String())
	}

	pp = trend.NewPivotPointWithMethod[float64](trend.PivotPointWoodie)
	if pp.String() != "PivotPoint(Woodie)" {
		t.Fatalf("expected PivotPoint(Woodie), got %s", pp.String())
	}

	pp = trend.NewPivotPointWithMethod[float64](trend.PivotPointCamarilla)
	if pp.String() != "PivotPoint(Camarilla)" {
		t.Fatalf("expected PivotPoint(Camarilla), got %s", pp.String())
	}

	pp = trend.NewPivotPointWithMethod[float64](trend.PivotPointFibonacci)
	if pp.String() != "PivotPoint(Fibonacci)" {
		t.Fatalf("expected PivotPoint(Fibonacci), got %s", pp.String())
	}
}
