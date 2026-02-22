// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestConnorsRsi(t *testing.T) {
	type Data struct {
		Close      float64
		ConnorsRsi float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/connors_rsi.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.ConnorsRsi })

	connorsRsi := momentum.NewConnorsRsi[float64]()
	actual := connorsRsi.Compute(closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, connorsRsi.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConnorsRsiString(t *testing.T) {
	connorsRsi := momentum.NewConnorsRsi[float64]()
	expected := "ConnorsRSI(3, 2, 100)"
	actual := connorsRsi.String()

	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
