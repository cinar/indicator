// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
)

func TestStochasticRsi(t *testing.T) {
	type Data struct {
		Close         float64
		StochasticRsi float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic_rsi.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closings := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.StochasticRsi })

	stochasticRsi := momentum.NewStochasticRsi[float64]()
	actual := stochasticRsi.Compute(closings)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, stochasticRsi.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
