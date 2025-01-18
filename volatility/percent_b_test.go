// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volatility_test

import (
	"fmt"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/volatility"
)

func ExamplePercentB() {
	// Closing prices
	closes := helper.SliceToChan([]float64{
		318.600006, 315.839996, 316.149994, 310.570007, 307.779999,
		305.820007, 305.98999, 306.390015, 311.450012, 312.329987,
		309.290009, 301.910004, 300, 300.029999, 302,
		307.820007, 302.690002, 306.48999, 305.549988, 303.429993,
	})

	// Initialize the %B indicator
	percentB := volatility.NewPercentB[float64]()

	// Compute %B
	result := percentB.Compute(closes)

	// Round digits
	result = helper.RoundDigits(result, 2)

	fmt.Println(helper.ChanToSlice(result))
	// Output: [0.3]
}

func TestPercentB(t *testing.T) {
	type Data struct {
		Close    float64
		PercentB float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/percent_b.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })
	expected := helper.Map(inputs[1], func(d *Data) float64 { return d.PercentB })

	percentB := volatility.NewPercentB[float64]()

	actual := percentB.Compute(closing)
	actual = helper.RoundDigits(actual, 2)

	expected = helper.Skip(expected, percentB.IdlePeriod())

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPercentBString(t *testing.T) {
	expected := "%B(10)"
	actual := volatility.NewPercentBWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
