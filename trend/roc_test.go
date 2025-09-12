package trend

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRocSimple(t *testing.T) {
	closing := helper.SliceToChan([]float64{2, 4, 6, 8, 9, 7})
	expected := helper.SliceToChan([]float64{0, 0, 0, (8 - 2) / 2, (9.0 - 4) / 4, (7.0 - 6) / 6})

	roc := NewRoc[float64]()
	roc.Period = 3

	actual := roc.Compute(closing)
	expected = helper.Skip(expected, roc.IdlePeriod())

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRocTestdata(t *testing.T) {
	type Data struct {
		Close float64
		Roc   float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/roc.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.Roc })

	roc := NewRoc[float64]()

	actual := roc.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, roc.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRocFallbackPeriod(t *testing.T) {
	roc := NewRocWithPeriod[float64](-17)

	if roc.Period != DefaultRocPeriod {
		t.Fatal("expected period to be fallback to default value")
	}
}

func TestRocToStringAndIdlePeriod(t *testing.T) {
	roc := NewRocWithPeriod[float64](0)
	if roc.IdlePeriod() != DefaultRocPeriod {
		t.Fatalf("unexpected IdlePeriod: %d", roc.IdlePeriod())
	}
	roc.Period = 3
	if s := roc.String(); s != "ROC(3)" {
		t.Fatalf("unexpected String(): %s", s)
	}
}
