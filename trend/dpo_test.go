package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestDpo(t *testing.T) {
	type DpoData struct {
		Close float64
		Dpo   float64
	}

	input, err := helper.ReadFromCsvFile[DpoData]("testdata/dpo.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *DpoData) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *DpoData) float64 { return d.Dpo })

	dpo := trend.NewDpo[float64]()
	actual := helper.RoundDigits(dpo.Compute(closing), 2)
	expected = helper.Skip(expected, dpo.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDpoString(t *testing.T) {
	expected := "DPO(10)"
	actual := trend.NewDpoWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
