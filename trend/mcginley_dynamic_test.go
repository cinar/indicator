// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestMcGinleyDynamic(t *testing.T) {
	type Data struct {
		Close    float64 `header:"Close"`
		Expected float64 `header:"Expected"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/mcginley_dynamic.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	ind := trend.NewMcGinleyDynamic[float64]()
	actuals := ind.Compute(closings)
	actuals = helper.RoundDigits(actuals, 6)

	inputs[1] = helper.Skip(inputs[1], ind.IdlePeriod())

	for data := range inputs[1] {
		actual, ok := <-actuals
		if !ok {
			t.Fatal("actuals channel closed early")
		}
		if actual != data.Expected {
			t.Fatalf("actual %v expected %v", actual, data.Expected)
		}
	}
}

func TestNewMcGinleyDynamic(t *testing.T) {
	ind := trend.NewMcGinleyDynamic[float64]()
	if ind.Period != 14 {
		t.Fatalf("expected 14 but got %d", ind.Period)
	}
}

func TestNewMcGinleyDynamicWithPeriod(t *testing.T) {
	ind := trend.NewMcGinleyDynamicWithPeriod[float64](20)
	if ind.Period != 20 {
		t.Fatalf("expected 20 but got %d", ind.Period)
	}
}

func TestMcGinleyDynamicString(t *testing.T) {
	ind := trend.NewMcGinleyDynamicWithPeriod[float64](14)
	if ind.String() != "MD(14)" {
		t.Fatalf("expected MD(14) but got %s", ind.String())
	}
}

func TestMcGinleyDynamicZero(t *testing.T) {
	closings := helper.SliceToChan([]float64{0, 10, 20})
	ind := trend.NewMcGinleyDynamicWithPeriod[float64](14)
	actuals := ind.Compute(closings)

	if <-actuals != 0 {
		t.Fatal("expected 0")
	}
	if <-actuals != 10 {
		t.Fatal("expected 10")
	}
	// (20-10) / (14 * (20/10)^4) = 10 / (14 * 16) = 10 / 224 = 0.044642857
	// 10 + 0.044642857 = 10.044642857
	actual := <-actuals
	expected := 10.044642857
	if math.Abs(float64(actual)-expected) > 0.000001 {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}
