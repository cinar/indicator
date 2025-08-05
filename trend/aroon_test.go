// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestAroon(t *testing.T) {
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}

	aroon := trend.NewAroon[float64]()

	input, err := helper.ReadFromCsvFile[AroonData]("testdata/aroon.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row *AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row *AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row *AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row *AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.Period-1)
	expectedDown = helper.Skip(expectedDown, aroon.Period-1)

	actualUp, actualDown := aroon.Compute(high, low)

	err = helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAroonCompareRust compares the result of the Go implementation with the result of a Rust implementation.
func TestAroonCompareRust(t *testing.T) {
	// compares to: https://github.com/sabinchitrakar/aroon-rs/blob/master/docs/aroon.md
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}
	var testData = []AroonData{
		{High: 82.15, Low: 81.29, Up: 0, Down: 0}, // None
		{High: 81.89, Low: 80.64, Up: 0, Down: 0}, // None
		{High: 83.03, Low: 81.31, Up: 0, Down: 0}, // None
		{High: 83.30, Low: 82.65, Up: 0, Down: 0}, // None
		{High: 83.85, Low: 83.07, Up: 0, Down: 0}, // None
		{High: 83.90, Low: 83.11, Up: 20.00, Down: 100.00},
		{High: 83.33, Low: 82.49, Up: 20.00, Down: 80.00},
		{High: 84.30, Low: 82.30, Up: 100.00, Down: 100.00},
		{High: 84.84, Low: 84.15, Up: 80.00, Down: 100.00},
		{High: 85.00, Low: 84.11, Up: 60.00, Down: 100.00},
		{High: 85.90, Low: 84.03, Up: 40.00, Down: 100.00},
		{High: 86.58, Low: 85.39, Up: 20.00, Down: 100.00},
		{High: 86.98, Low: 85.76, Up: 60.00, Down: 100.00},
		{High: 88.00, Low: 87.17, Up: 40.00, Down: 100.00},
		{High: 87.87, Low: 87.01, Up: 20.00, Down: 80.00},
	}

	input := helper.SliceToChan[AroonData](testData)
	aroon := trend.NewAroon[float64]()

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.Period-1)
	expectedDown = helper.Skip(expectedDown, aroon.Period-1)

	actualUp, actualDown := aroon.Compute(high, low)

	err := helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}

/* https://github.com/francescosisini/TTR/blob/27ea28698295f56447fcc87ad515140bdb35c8a5/R/aroon.R#L27
'Aroon up (down) is the elapsed time, expressed as a percentage, between today
#'and the highest (lowest) price in the last \code{n} periods.  If today's
#'price is a new high (low) Aroon up (down) will be 100. Each subsequent period
#'without another new high (low) causes Aroon up (down) to decrease by (1 /
#'\code{n}) x 100.
*/

func TestAroonValidateDefinition(t *testing.T) {
	// compares to: https://github.com/sabinchitrakar/aroon-rs/blob/master/docs/aroon.md
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}
	var testData = []AroonData{
		{High: 11, Low: 9, Up: 0, Down: 0},    // None
		{High: 10, Low: 9, Up: 0, Down: 0},    // None
		{High: 9, Low: 9, Up: 100, Down: 100}, // None
		{High: 9, Low: 9, Up: 67, Down: 67},   // None
		{High: 9, Low: 9, Up: 33, Down: 33},   // None
	}

	input := helper.SliceToChan[AroonData](testData)
	aroon := &trend.Aroon[float64]{
		Period: 3,
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.Period-1)
	expectedDown = helper.Skip(expectedDown, aroon.Period-1)

	actualUp, actualDown := aroon.Compute(high, low)

	err := helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}
